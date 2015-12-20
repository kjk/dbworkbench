/* jshint -W097,-W117 */
'use strict';

const React = require('react');

class AlertBar extends React.Component {
  render() {
    return <div id="alert-bar">{this.props.errorMessage}</div>;
  }
}

module.exports = AlertBar;