export namespace main {
	
	export class Query {
	    ID?: number;
	    Title: string;
	    Query: string;
	    Description: string;
	    CreatedAt?: string;
	    UpdatedAt?: string;
	    DeletedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new Query(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Title = source["Title"];
	        this.Query = source["Query"];
	        this.Description = source["Description"];
	        this.CreatedAt = source["CreatedAt"];
	        this.UpdatedAt = source["UpdatedAt"];
	        this.DeletedAt = source["DeletedAt"];
	    }
	}
	export class Variable {
	    Field: string;
	    Value: string;
	    Position: number;
	
	    static createFrom(source: any = {}) {
	        return new Variable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Field = source["Field"];
	        this.Value = source["Value"];
	        this.Position = source["Position"];
	    }
	}

}

