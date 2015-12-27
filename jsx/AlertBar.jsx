import React from 'react';

export default class AlertBar extends React.Component {
  render() {
    return <div id="alert-bar">{this.props.errorMessage}</div>;
  }
}
