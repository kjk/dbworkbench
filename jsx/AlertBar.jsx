import React from 'react';

const AlertBar = (props) => {
  return <div id="alert-bar">
           { props.errorMessage }
         </div>;
};

export default AlertBar;

