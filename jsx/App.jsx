/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');
var ConnectionWindow = require('./ConnectionWindow.jsx');
var TopNav = require('./TopNav.jsx');
var Sidebar = require('./Sidebar.jsx');
var Input = require('./Input.jsx');
var Output = require('./Output.jsx');

var App = React.createClass({
  getInitialState: function() {
    return {
      connectionId: -1,
      connected: false,
      results: null,
    };
  },

  handleDidConnect: function(connectionStr, connectionId) {
    console.log("App.handleDidConnect: ", connectionStr, connectionId);
    this.setState({
      connected: true,
      connectionId: connectionId
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
          <TopNav />
          <Sidebar onGotResults={this.handleGotResults}/>
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
