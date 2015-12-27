import React from 'react';

let AlertBar = (props) => {
  return <div id="alert-bar">{props.errorMessage}</div>;
};

module.exports = AlertBar;