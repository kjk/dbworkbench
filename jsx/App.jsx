/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');
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
      connectionId: -1,
      connected: false,
      databaseName: "",
      selectedView: view.SQLQuery,
      results: null,
    };
  },

  handleDidConnect: function(connectionStr, connectionId, databaseName) {
    console.log("App.handleDidConnect: ", connectionStr, connectionId, databaseName);
    this.setState({
      connected: true,
      connectionId: connectionId,
      databaseName: databaseName
    });
  },

  handleGotResults: function(results) {
    console.log("handleGotResults: ", results);
    this.setState({
      results: results
    });
  },

  renderInput: function() {
    return <Input />;
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
  },

  componentDidUnmount: function() {
    //action.cancelOnViewSelected(this.handleViewSelected);
  },

  render: function() {
    var results = this.state.results;
    var input;
    if (!results) {
      input = this.renderInput();
    }

    if (!this.state.connected) {
      return <ConnectionWindow onDidConnect={this.handleDidConnect} />;
    } else {
      return (
        <div>
          <TopNav view={this.state.selectedView}/>
          <Sidebar
            onGotResults={this.handleGotResults}
            databaseName={this.state.databaseName}
          />
          <div id="body">
            {input}
            <Output results={results}/>
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
