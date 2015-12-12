/* jshint -W097,-W117 */
'use strict';

var React = require('react');

// var action = require('./action.js');
// var api = require('./api.js');
var DbNav = require('./DbNav.jsx');
var Input = require('./Input.jsx');
var Output = require('./Output.jsx');
var view = require('./view.js');

class MainContainer extends React.Component {
  renderInput(tooLong) {
    if (this.props.selectedView === view.SQLQuery) {
      return <Input tooLong={tooLong} />;
    }
  }

  render() {
    // when showing sql query, results are below editor window
    var notFull = (this.props.selectedView === view.SQLQuery);

    var divStyle = {
      left: this.props.dragBarPosition + 'px',
    }

    var results = this.props.results
    if (results && this.props.results.rows.length > 100) {
      // It's only showed when +100. We could make this default.
      var tooLong = "Showing 100 out of " + results.rows.length + " rows."
      results.rows = results.rows.slice(0, 100);
    }

    return (
      <div id="body" style={divStyle}>
          <DbNav view={this.props.selectedView}/>
          {this.renderInput(tooLong)}
          <Output
            selectedView={this.props.selectedView}
            results={results}
            notFull={notFull}/>
      </div>
    );
  }
}

module.exports = MainContainer;