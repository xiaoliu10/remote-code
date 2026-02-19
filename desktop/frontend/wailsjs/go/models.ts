export namespace main {
	
	export class Config {
	    backendPort: number;
	    frontendPort: number;
	    jwtSecret: string;
	    adminPassword: string;
	    allowedDir: string;
	    frpEnabled: boolean;
	    frpServer: string;
	    frpPort: number;
	    frpToken: string;
	    setupCompleted: boolean;
	    platform: string;
	    isLocalSupported: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.backendPort = source["backendPort"];
	        this.frontendPort = source["frontendPort"];
	        this.jwtSecret = source["jwtSecret"];
	        this.adminPassword = source["adminPassword"];
	        this.allowedDir = source["allowedDir"];
	        this.frpEnabled = source["frpEnabled"];
	        this.frpServer = source["frpServer"];
	        this.frpPort = source["frpPort"];
	        this.frpToken = source["frpToken"];
	        this.setupCompleted = source["setupCompleted"];
	        this.platform = source["platform"];
	        this.isLocalSupported = source["isLocalSupported"];
	    }
	}

}

