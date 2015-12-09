/* jshint -W097,-W117 */
'use strict';

var React = require('react');

class AlertBar extends React.Component {
  render() {
    return <div id="note">{this.props.errorMessage}</div>;
  }
}

module.exports = AlertBar;