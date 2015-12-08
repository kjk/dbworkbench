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
  renderInput() {
    if (this.props.selectedView === view.SQLQuery) {
      return <Input />;
    }
  }

  render() {
    // when showing sql query, results are below editor window
    var notFull = (this.props.selectedView === view.SQLQuery);

    var divStyle = {
      left: this.props.dragBarPosition + 'px',
    }

    return (
      <div id="body" style={divStyle}>
          <DbNav view={this.props.selectedView}/>
          {this.renderInput()}
          <Output
            selectedView={this.props.selectedView}
            results={this.props.results}
            notFull={notFull}/>
      </div>
    );
  }
}

module.exports = MainContainer;