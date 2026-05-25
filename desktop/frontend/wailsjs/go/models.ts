export namespace main {
	
	export class Status {
	    agent_status: string;
	    device_status: string;
	    protection_status: string;
	    backend_url: string;
	    server_status: string;
	    protection_enabled: boolean;
	    pairing_status: string;
	    device_code: string;
	    pairing_expires_at: string;
	    mqtt_status: string;
	    last_error: string;
	    last_event: string;
	    last_event_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.agent_status = source["agent_status"];
	        this.device_status = source["device_status"];
	        this.protection_status = source["protection_status"];
	        this.backend_url = source["backend_url"];
	        this.server_status = source["server_status"];
	        this.protection_enabled = source["protection_enabled"];
	        this.pairing_status = source["pairing_status"];
	        this.device_code = source["device_code"];
	        this.pairing_expires_at = source["pairing_expires_at"];
	        this.mqtt_status = source["mqtt_status"];
	        this.last_error = source["last_error"];
	        this.last_event = source["last_event"];
	        this.last_event_at = source["last_event_at"];
	    }
	}

}

