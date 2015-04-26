/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');
var ConnectionWindow = require('./ConnectionWindow.jsx');

var App = React.createClass({
  getInitialState: function() {
    return {
      connectionId: -1,
      connected: false,
    };
  },

  handleDidConnect: function(connectionStr, connectionId) {
    console.log("App.handleDidConnect: ", connectionStr, connectionId);
  },

  render: function() {
    if (this.state.connectionId === -1) {
      return <ConnectionWindow onDidConnect={this.handleDidConnect} />;
    } else {
      return <div>This is a start</div>;
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
