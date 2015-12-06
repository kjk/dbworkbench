/* jshint -W097,-W117 */
'use strict';

var Content = "Content";
var Structure = "Structure";
var Indexes = "Indexes";
var SQLQuery = "SQL Query";
var History = "History";
var Activity = "Activity";
var Connection = "Connection";

var AllViews = [
  SQLQuery,
  Content,
  Structure,
  Indexes,
  History,
  Activity,
  Connection
];

module.exports = {
  SQLQuery: SQLQuery,
  Content: Content,
  Structure: Structure,
  Indexes: Indexes,
  History: History,
  Activity: Activity,
  Connection: Connection,
  AllViews: AllViews
};
