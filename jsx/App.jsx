/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');
var ConnectionWindow = require('./ConnectionWindow.jsx');

var App = React.createClass({
  getInitialState: function() {
    return {
      connectionId: null,
    }
  },

  render: function() {
    if (this.state.connectionId === null) {
      return <ConnectionWindow />;
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
