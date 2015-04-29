/* jshint -W097,-W117 */
'use strict';

/*
function buildTable(results, sortColumn, sortOrder) {
  results.columns.forEach(function(col) {
      if (col === sortColumn) {
          cols += "<th data='" + col + "'" + "data-sort-order=" + sortOrder + ">" + col + "&nbsp;" + sortArrow(sortOrder) + "</th>";
      } else {
          cols += "<th data='" + col + "'>" + col + "</th>";
      }
  });
}
*/

var Output = React.createClass({
  renderHeader: function(columns, sortColumn, sortOrder) {
    var i = 0;
    var children = columns.map(function(col) {
      // TODO: use sortColumn and sortOrder
      i = i + 1;
      return (
        <th key={i} data={col}>{col}</th>
      );
    });

    return (
      <thead><tr>{children}</tr></thead>
    );
  },

  renderRow: function(row, key) {
    var i = 0;
    var children = row.map(function(col) {
      i = i + 1;
      return (
        <td key={i}><div>{col}</div></td>
      );
    });
    return (
      <tr key={key}>{children}</tr>
    );
  },

  renderRows: function(rows) {
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
  },

  renderResults: function(results) {
    var header = this.renderHeader(results.columns);
    var rows = this.renderRows(results.rows);
    return (
      <table id="results" className="table" data-mode="browse">
        {header}
        {rows}
      </table>
    );
  },

  renderNoResults: function() {
    return (
      <table id="results" className="table empty">
        <tbody>
          <tr><td>No records found</td></tr>
        </tbody>
      </table>
    );
  },

  renderError: function(errorMsg) {
    return (
      <table id="results" className="table empty">
        <tbody>
          <tr><td>ERROR: {errorMsg}</td></tr>
        </tbody>
      </table>
    );
  },

  render: function() {
    var clsOutput;
    var results;
    if (!this.props.results) {
      results = this.renderNoResults();
    } else {
      if (this.props.results.error) {
        results = this.renderError(this.props.results.error);
      } else {
        clsOutput = "full";
        results = this.renderResults(this.props.results);
      }
    }

    return (
      <div id="output" className={clsOutput}>
        <div className="wrapper">
            {results}
        </div>
      </div>
    );
  }
});

module.exports = Output;
