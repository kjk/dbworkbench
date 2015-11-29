/* jshint -W097,-W117 */
'use strict';

var AlertBar = React.createClass({
  render: function(){
    return <div id="note">{this.props.errorMessage}</div>;
  }
});

module.exports = AlertBar;