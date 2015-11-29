/* jshint -W097,-W117 */
'use strict';

var AlertBar = React.createClass({
  render: function(){
    return <div id="note">{this.props.errorMessage} <span><strong>Close</strong></span></div>;
  }
});

module.exports = AlertBar;