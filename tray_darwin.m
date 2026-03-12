// tray_darwin.m — macOS NSStatusItem (menu-bar icon) for GoTube.
// CGo compiles ObjC WITHOUT ARC → explicit retain/release (MRC).
//
// Standard macOS menu-bar behavior: click icon → dropdown menu.

#import <Cocoa/Cocoa.h>
#include <stdio.h>

// Go-exported callbacks (dispatch to goroutines on Go side)
extern void goTrayToggle(void);
extern void goTrayQuit(void);

// ─── Forward declarations & statics ─────────────────────────────────────────

@class GoTrayDelegate;

static NSStatusItem   *_statusItem = nil;
static NSMenu         *_trayMenu   = nil;
static GoTrayDelegate *_delegate   = nil;

// ─── Delegate – handles menu actions ────────────────────────────────────────

@interface GoTrayDelegate : NSObject <NSMenuDelegate>
- (void)openAction:(id)sender;
- (void)quitAction:(id)sender;
@end

@implementation GoTrayDelegate

- (void)openAction:(id)sender {
    goTrayToggle();
}

- (void)quitAction:(id)sender {
    goTrayQuit();
}

// Highlight the button while the menu is open (standard macOS behavior)
- (void)menuWillOpen:(NSMenu *)menu {
    [[_statusItem button] setHighlighted:YES];
}

- (void)menuDidClose:(NSMenu *)menu {
    [[_statusItem button] setHighlighted:NO];
}

@end

// ─── doSetupTray — runs on the main thread ──────────────────────────────────

static void doSetupTray(NSData *icon1xData, NSData *icon2xData) {
    fprintf(stderr, "[tray] doSetupTray: main=%d\n", (int)[NSThread isMainThread]);

    // 1. Delegate
    _delegate = [[GoTrayDelegate alloc] init];

    // 2. Build menu — standard macOS dropdown
    _trayMenu = [[NSMenu alloc] initWithTitle:@"GoTube"];
    [_trayMenu setDelegate:_delegate];

    // "Open GoTube" with hotkey hint
    NSMenuItem *openItem = [[NSMenuItem alloc]
        initWithTitle:@"Open GoTube"
               action:@selector(openAction:)
        keyEquivalent:@""];
    [openItem setTarget:_delegate];
    [_trayMenu addItem:openItem];
    [openItem release];

    [_trayMenu addItem:[NSMenuItem separatorItem]];

    // "Quit GoTube" with ⌘Q equivalent
    NSMenuItem *quitItem = [[NSMenuItem alloc]
        initWithTitle:@"Quit GoTube"
               action:@selector(quitAction:)
        keyEquivalent:@"q"];
    [quitItem setTarget:_delegate];
    [_trayMenu addItem:quitItem];
    [quitItem release];

    // 3. Create status item with fixed square size
    _statusItem = [[[NSStatusBar systemStatusBar]
                     statusItemWithLength:NSSquareStatusItemLength] retain];

    // 4. Build template icon from 1x + 2x PNG data
    NSImage *icon = nil;
    if (icon1xData && [icon1xData length] > 0) {
        icon = [[NSImage alloc] initWithData:icon1xData];
        if (icon) {
            [icon setSize:NSMakeSize(18, 18)];

            // Add @2x representation for Retina
            if (icon2xData && [icon2xData length] > 0) {
                NSBitmapImageRep *rep2x = [[NSBitmapImageRep alloc] initWithData:icon2xData];
                if (rep2x) {
                    [rep2x setSize:NSMakeSize(18, 18)];
                    [icon addRepresentation:rep2x];
                    [rep2x release];
                }
            }

            // Template mode: macOS auto-tints for dark/light mode
            [icon setTemplate:YES];
        }
    }

    // 5. Configure button
    NSStatusBarButton *btn = [_statusItem button];
    if (btn) {
        if (icon) {
            [btn setImage:icon];
            [btn setImagePosition:NSImageOnly];
        } else {
            [btn setTitle:@"GT"];
        }
        [btn setToolTip:@"GoTube Downloader — ⌃⇧Y"];
    }
    if (icon) [icon release];

    // 6. Attach menu — standard macOS: ANY click on the icon opens this menu
    [_statusItem setMenu:_trayMenu];

    fprintf(stderr, "[tray] setup complete — icon in menu bar\n");
}

// ─── C entry points called from Go via CGo ──────────────────────────────────

void nativeSetupTray(const void *icon1x, int icon1xLen,
                     const void *icon2x, int icon2xLen) {
    fprintf(stderr, "[tray] nativeSetupTray: 1x=%d bytes, 2x=%d bytes\n",
            icon1xLen, icon2xLen);

    // Hide from Dock IMMEDIATELY
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    });

    // Copy icon data (retain for delayed use under MRC)
    NSData *data1x = [[NSData dataWithBytes:icon1x length:icon1xLen] retain];
    NSData *data2x = [[NSData dataWithBytes:icon2x length:icon2xLen] retain];

    // Setup after Wails finishes its own window init
    dispatch_after(
        dispatch_time(DISPATCH_TIME_NOW, (int64_t)(1.5 * NSEC_PER_SEC)),
        dispatch_get_main_queue(),
        ^{
            doSetupTray(data1x, data2x);
            [data1x release];
            [data2x release];
        });
}

void nativeTeardownTray(void) {
    dispatch_async(dispatch_get_main_queue(), ^{
        if (_statusItem) {
            [[NSStatusBar systemStatusBar] removeStatusItem:_statusItem];
            [_statusItem release];
            _statusItem = nil;
        }
        if (_trayMenu)  { [_trayMenu release];  _trayMenu = nil; }
        if (_delegate)  { [_delegate release];  _delegate = nil; }
    });
}
