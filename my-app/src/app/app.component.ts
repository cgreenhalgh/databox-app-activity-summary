import { Component } from '@angular/core';
import { OnInit } from '@angular/core';

import { DataService } from './data.service';
import { DataType } from './types';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  loading = true;
  dataTypes: DataType[];
  constructor(private dataService: DataService) {}
  ngOnInit(): void {
	  this.dataService.getDataTypes().then(dts => { this.dataTypes = dts; this.loading = false });
  }
}
