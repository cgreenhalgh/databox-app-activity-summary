import { Component } from '@angular/core';
import { OnInit } from '@angular/core';

import { DataService } from './data.service';
import { DataType } from './types';

import * as vega from 'vega';
import * as vegaLite from 'vega-lite';

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
  
    // test example
    let vlSpec = getExVegaLiteSpec();
    let theSpec = vegaLite.compile(vlSpec).spec;
   
    const container = document.querySelector("#view");
    const runtime = vega.parse(theSpec, {});
    const view = new vega.View(runtime, { logLevel: vega.Warn });
    const initializedView = view.initialize(container);
    view.run();
  }
}
function getExVegaLiteSpec(): any {
  return {
    "$schema": "https://vega.github.io/schema/vega-lite/v2.json",
    "description": "A simple bar chart with embedded data.",
    "width": 70,
    "data": { 
      "url": "../api/get/strava_activity/"
    },
    "mark": "bar",
    "transform": [
      {"calculate": "datum.data.distance*0.001", "as": "distance"}
      //,{"calculate": "floor(datum.timestamp / (24*60*60*1000))*(24*60*60*1000)", "as": "date"},
      //,{"filter": "datum.timestamp > now()-7*24*60*60*1000"}
      ,{"calculate": "floor((now()-datum.timestamp)/(7*24*60*60*1000))", "as": "weeksAgo"}
    ],
    "encoding": {
      "column": {
          "field": "weeksAgo", "type": "ordinal"
      },
      "x": {
        //"field": "date",
        "field": "timestamp",
        "type": "temporal"
        ,"timeUnit": "day"
        ,"axis": {"title": "day", "grid": false}
      },
      "y": {
        "field": "distance",
        "type": "quantitative",
        "aggregate": "sum",
        "axis":{
          "title":"Total distance (km)"
        }
      }
      ,"color": {
        "field": "weeksAgo", "type": "nominal"
      }
    },
    "config": {
      "view": {"stroke": "transparent"},
      "axis": {"domainWidth": 1}
    }
  }
}
