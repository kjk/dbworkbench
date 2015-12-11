/* jshint -W097,-W117 */
'use strict';

var React = require('react');
var _ = require('underscore');

var view = require('./view.js');

class Output extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleRowClick = this.handleRowClick.bind(this);

    this.state = {
      clickedRowKey: -1,
      rowStyle: {},
    };
  }

  handleRowClick(key, e) {
    console.log("Enlarging ", key)
    var enlargeStyle = {
      maxWidth: '350px',
      maxHeight: '100%',
      overflow: 'hidden',
      textOverflow: 'ellipsis',
      whiteSpace: 'normal',
    }

    if (_.isEqual(enlargeStyle, this.state.rowStyle) && key == this.state.clickedRowKey) {
      console.log("Shrinking")
      enlargeStyle = {}
    }

    this.setState({
      clickedRowKey: key,
      rowStyle: enlargeStyle
    });
  }

  renderHeader(columns, sortColumn, sortOrder) {
    var i = 0;
    if (!columns) {
      columns = [];
    }
    var children = columns.map(function(col) {
      // TODO: use sortColumn and sortOrder)
      i = i + 1;
      return (
        <th key={i} data={col}>{col}</th>
      );
    });

    return (
      <thead><tr>{children}</tr></thead>
    );
  }

  renderRow(row, key) {
    var style = {}
    if (this.state.clickedRowKey == key) {
      style = this.state.rowStyle
    }

    var i = 0;
    var children = row.map(function(col) {
      i = i + 1;
      return (
        <td key={i}><div style={style}>{col}</div></td>
      );
    });

    return (
      <tr key={key} onClick={this.handleRowClick.bind(this,key)}>{children}</tr>
    );
  }

  renderRows(rows) {
    if (!rows) {
      return;
    }

    var self = this;
    var i = 0;
    var children = rows.map(function(row) {
      i = i + 1;
      return self.renderRow(row, i);
    });

    return (
      <tbody>{children}</tbody>
    );
  }

  renderResults(results) {
    var header = this.renderHeader(results.columns);
    var rows = this.renderRows(results.rows);
    return (
      <table id="results" className="table" data-mode="browse">
        {header}
        {rows}
      </table>
    );
  }

  renderNoResults() {
    return (
      <table id="results" className="table empty no-crop">
        <tbody>
          <tr><td>No records found</td></tr>
        </tbody>
      </table>
    );
  }

  renderError(errorMsg) {
    return (
      <table id="results" className="table empty">
        <tbody>
          <tr><td>ERROR: {errorMsg}</td></tr>
        </tbody>
      </table>
    );
  }

  render() {
    var clsOutput, children;
    var results = this.props.results;
    if (!results) {
      children = this.renderNoResults();
    } else {
      if (results.error) {
        children = this.renderError(results.error);
      } else if (!results.rows || results.rows.length === 0) {
        children = this.renderNoResults();
        clsOutput = "full";
      } else {
        clsOutput = "full";
        children = this.renderResults(results);
      }
    }
    if (this.props.notFull) {
      clsOutput = "";
    }

    if (view.SQLQuery != this.props.selectedView) {
      clsOutput = "full";
    }

    return (
      <div id="output" className={clsOutput}>
        <div className="wrapper">
            {children}
        </div>
      </div>
    );
  }
}

module.exports = Output;
