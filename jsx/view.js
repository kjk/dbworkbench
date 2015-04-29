/* jshint -W097,-W117 */
'use strict';

var Content = 0;
var Structure = 1;
var Indexes = 2;
var SQLQuery = 3;
var History = 4;
var Activity = 5;
var Connection = 6;

var Names = [
  "Content",
  "Structure",
  "Indexes",
  "SQL Query",
  "History",
  "Activity",
  "Connection"
];

var AllViews = [
  Content,
  Structure,
  Indexes,
  SQLQuery,
  History,
  Activity,
  Connection
];

exports.Content = Content;
exports.Structure = Structure;
exports.Indexes = Indexes;
exports.SQLQuery = SQLQuery;
exports.History = History;
exports.Activity = Activity;
exports.Connection = Connection;
exports.Names = Names;
exports.AllViews = AllViews;
