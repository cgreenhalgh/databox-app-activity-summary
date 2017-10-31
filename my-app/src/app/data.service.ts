import { Injectable } from '@angular/core';
import { Headers, Http } from '@angular/http';

import 'rxjs/add/operator/toPromise';

import { DataType } from './types';

@Injectable()
export class DataService {
	private dataTypesUrl = '../api/dataTypes';  // URL to web api

	constructor(private http: Http) {
		console.log('creating DataService');
	}
	getDataTypes(): Promise<DataType[]> {
		//const test : DataType[] = [ { id:"1", title:"something", available: true }]
		//return Promise.resolve(test);
		return this.http.get(this.dataTypesUrl)
			.toPromise()
			.then(response => response.json() as DataType[])
			.catch(this.handleError);
	}
		 
	private handleError(error: any): Promise<any> {
		console.error('An error occurred', error); // for demo purposes only
		return Promise.reject(error.message || error);
	}
}