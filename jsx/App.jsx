/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');
var api = require('./api.js');
var action = require('./action.js');
var view = require('./view.js');

var ConnectionWindow = require('./ConnectionWindow.jsx');
var TopNav = require('./TopNav.jsx').TopNav;
var ViewSQLQuery = require('./TopNav.jsx').ViewSQLQuery;
var Sidebar = require('./Sidebar.jsx');
var Input = require('./Input.jsx');
var Output = require('./Output.jsx');

var App = React.createClass({
  getInitialState: function() {
    return {
      selectedView: view.SQLQuery,
      connectionId: -1,
      connected: false,
      databaseName: "",
      tables: null,
      selectedTable: "",
      selectedTableInfo: null,
      results: null,
    };
  },

  handleDidConnect: function(connectionStr, connectionId, databaseName) {
    this.setState({
      connected: true,
      connectionId: connectionId,
      databaseName: databaseName
    });
    var self = this;
    api.getTables(function(data) {
      self.setState({
        tables: data,
      });
    });
  },

  renderInput: function() {
    if (this.state.results) {
      return <Input />;
    }
  },

  handleTableSelected: function(table) {
    this.setState({
      selectedTable: table
    });

    var self = this;
    api.call("get", "/tables/" + table + "/info", {}, function(data) {
      console.log("handleSelectTable: tableInfo: ", data);
      self.setState({
        selectedTableInfo: data,
      });
    });

    var sortColumn = null;
    var sortOrder = null;
    var params = { limit: 100, sort_column: sortColumn, sort_order: sortOrder };
    api.getTableRows(table, params, function(data) {
      console.log("handleSelectTable: got table rows: ", data);
      self.setState({
        results: data
      });
      action.viewSelected(view.Content);
    });
  },

  handleViewSelected: function(view) {
    console.log("handleViewSelected: ", view);
    // TODO: load the right data as results
    this.setState({
      selectedView: view
    });
  },

  componentDidMount: function() {
    action.onViewSelected(this.handleViewSelected);
    action.onTableSelected(this.handleTableSelected);
  },

  componentDidUnmount: function() {
    //action.cancelOnViewSelected(this.handleViewSelected);
  },

  render: function() {
    if (!this.state.connected) {
    return <ConnectionWindow onDidConnect={this.handleDidConnect} />;
  } else {
    return (
      <div>
        <TopNav view={this.state.selectedView}/>
        <Sidebar
          tables={this.state.tables}
          selectedTable={this.state.selectedTable}
          selectedTableInfo={this.state.selectedTableInfo}
          databaseName={this.state.databaseName}
        />
        <div id="body">
          {this.renderInput()}
          <Output results={this.state.results}/>
        </div>
      </div>
    );
  }
  }
});

function appStart() {
  React.render(
    <App/>,
    document.getElementById('main')
  );
}

window.appStart = appStart;
