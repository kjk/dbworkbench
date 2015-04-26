/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');
var ConnectionWindow = require('./ConnectionWindow.jsx');
var TopNav = require('./Main.jsx').TopNav;
var Sidebar = require('./Main.jsx').Sidebar;
var Body = require('./Main.jsx').Body;

var App = React.createClass({
  getInitialState: function() {
    return {
      connectionId: -1,
      connected: false,
    };
  },

  handleDidConnect: function(connectionStr, connectionId) {
    console.log("App.handleDidConnect: ", connectionStr, connectionId);
    this.setState({
      connected: true,
      connectionId: connectionId
    });
  },

  renderMain: function() {
    return (
      <div>
        <TopNav />
        <Sidebar />
        <Body />
      </div>
    );
  },

  render: function() {
    if (!this.state.connected) {
      return <ConnectionWindow onDidConnect={this.handleDidConnect} />;
    } else {
      return this.renderMain();
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
