/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');

var App = React.createClass({
  render: function() {
    return <div>This is a start</div>;
  }
});

function appStart() {
  React.render(
    <App/>,
    document.getElementById('root')
  );
}

window.appStart = appStart;
