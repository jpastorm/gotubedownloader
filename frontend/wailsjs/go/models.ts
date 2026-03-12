export namespace bootstrap {
	
	export class BootstrapStatus {
	    ytDlpReady: boolean;
	    ffmpegReady: boolean;
	    ytDlpPath: string;
	    ffmpegPath: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new BootstrapStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ytDlpReady = source["ytDlpReady"];
	        this.ffmpegReady = source["ffmpegReady"];
	        this.ytDlpPath = source["ytDlpPath"];
	        this.ffmpegPath = source["ffmpegPath"];
	        this.error = source["error"];
	    }
	}

}

export namespace config {
	
	export class Config {
	    downloadDir: string;
	    maxConcurrent: number;
	    fragmentConcurrent: number;
	    defaultMode: string;
	    subtitleLang: string;
	    speedLimit: string;
	    proxy: string;
	    skipDuplicates: boolean;
	    continueDownloads: boolean;
	    embedMetadata: boolean;
	    downloadSubtitles: boolean;
	    embedSubtitles: boolean;
	    notifyOnComplete: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.downloadDir = source["downloadDir"];
	        this.maxConcurrent = source["maxConcurrent"];
	        this.fragmentConcurrent = source["fragmentConcurrent"];
	        this.defaultMode = source["defaultMode"];
	        this.subtitleLang = source["subtitleLang"];
	        this.speedLimit = source["speedLimit"];
	        this.proxy = source["proxy"];
	        this.skipDuplicates = source["skipDuplicates"];
	        this.continueDownloads = source["continueDownloads"];
	        this.embedMetadata = source["embedMetadata"];
	        this.downloadSubtitles = source["downloadSubtitles"];
	        this.embedSubtitles = source["embedSubtitles"];
	        this.notifyOnComplete = source["notifyOnComplete"];
	    }
	}

}

export namespace downloader {
	
	export class DownloadOptions {
	    url: string;
	    mode: string;
	    quality: string;
	    outputDir: string;
	    subtitles: boolean;
	    subtitleLang: string;
	    embedSubs: boolean;
	    embedMetadata: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DownloadOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.mode = source["mode"];
	        this.quality = source["quality"];
	        this.outputDir = source["outputDir"];
	        this.subtitles = source["subtitles"];
	        this.subtitleLang = source["subtitleLang"];
	        this.embedSubs = source["embedSubs"];
	        this.embedMetadata = source["embedMetadata"];
	    }
	}
	export class QueueItem {
	    id: string;
	    url: string;
	    title: string;
	    channel: string;
	    thumbnail: string;
	    duration: string;
	    status: string;
	    progress: number;
	    speed: string;
	    eta: string;
	    downloaded: string;
	    total: string;
	    error: string;
	    mode: string;
	    outputDir: string;
	    filePath: string;
	    logLines: string[];
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    completedAt?: any;
	    options: DownloadOptions;
	
	    static createFrom(source: any = {}) {
	        return new QueueItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.title = source["title"];
	        this.channel = source["channel"];
	        this.thumbnail = source["thumbnail"];
	        this.duration = source["duration"];
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.speed = source["speed"];
	        this.eta = source["eta"];
	        this.downloaded = source["downloaded"];
	        this.total = source["total"];
	        this.error = source["error"];
	        this.mode = source["mode"];
	        this.outputDir = source["outputDir"];
	        this.filePath = source["filePath"];
	        this.logLines = source["logLines"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
	        this.options = this.convertValues(source["options"], DownloadOptions);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VideoInfo {
	    id: string;
	    url: string;
	    title: string;
	    channel: string;
	    duration: number;
	    durationStr: string;
	    thumbnail: string;
	    type: string;
	    videoCount: number;
	    entries?: VideoInfo[];
	    isLive: boolean;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new VideoInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.title = source["title"];
	        this.channel = source["channel"];
	        this.duration = source["duration"];
	        this.durationStr = source["durationStr"];
	        this.thumbnail = source["thumbnail"];
	        this.type = source["type"];
	        this.videoCount = source["videoCount"];
	        this.entries = this.convertValues(source["entries"], VideoInfo);
	        this.isLive = source["isLive"];
	        this.description = source["description"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace history {
	
	export class HistoryEntry {
	    id: string;
	    url: string;
	    title: string;
	    channel: string;
	    duration: string;
	    thumbnail: string;
	    filePath: string;
	    mode: string;
	    quality: string;
	    status: string;
	    fileSize: number;
	    // Go type: time
	    downloadedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new HistoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.title = source["title"];
	        this.channel = source["channel"];
	        this.duration = source["duration"];
	        this.thumbnail = source["thumbnail"];
	        this.filePath = source["filePath"];
	        this.mode = source["mode"];
	        this.quality = source["quality"];
	        this.status = source["status"];
	        this.fileSize = source["fileSize"];
	        this.downloadedAt = this.convertValues(source["downloadedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

