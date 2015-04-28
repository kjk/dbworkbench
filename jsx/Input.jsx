/* jshint -W097,-W117 */
'use strict';

var Input = React.createClass({
  render: function() {
    return (
      <div id="input">
        <div className="wrapper">
          <div id="custom_query"></div>
          <div className="actions">
            <input type="button" id="run" value="Run Query" className="btn btn-sm btn-primary" />
            <input type="button" id="explain" value="Explain Query" className="btn btn-sm btn-default" />
            <input type="button" id="csv" value="Download CSV" className="btn btn-sm btn-default" />

            <div id="query_progress">Please wait, query is executing...</div>
          </div>
        </div>
      </div>
    );
  }
});

module.exports = Input;
