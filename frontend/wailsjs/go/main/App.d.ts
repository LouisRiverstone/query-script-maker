// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {sql} from '../models';

export function CreateOrUpdateDatabaseConnection(arg1:main.DatabaseConnection):Promise<main.DatabaseConnection>;

export function CreateSQLFile(arg1:string):Promise<string>;

export function DeleteQuery(arg1:number):Promise<void>;

export function ExportDatabaseFile():Promise<void>;

export function GetDatabaseConnection():Promise<main.DatabaseConnection>;

export function GetQueriesList(arg1:boolean):Promise<Array<main.Query>>;

export function ImportDatabaseFile():Promise<void>;

export function InsertQueryInDatabase(arg1:main.Query):Promise<void>;

export function MakeBindedSQL(arg1:string,arg2:Array<{[key: string]: any}>,arg3:Array<main.Variable>):Promise<string>;

export function ReadXLSXFile():Promise<string>;

export function TestQueryInDatabase(arg1:string):Promise<sql.Rows>;

export function UpdateQuery(arg1:number,arg2:main.Query):Promise<void>;
