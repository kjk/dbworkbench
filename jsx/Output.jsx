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
    if (!columns) {
      columns = [];
    }
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
      <table id="results" className="table empty no-crop">
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
    var clsOutput, children;
    var results = this.props.results;
    if (!results || !results.rows || results.rows.length === 0) {
      children = this.renderNoResults();
      clsOutput = "full";
    } else {
      if (results.error) {
        children = this.renderError(results.error);
      } else {
        clsOutput = "full";
        children = this.renderResults(results);
      }
    }

    return (
      <div id="output" className={clsOutput}>
        <div className="wrapper">
            {children}
        </div>
      </div>
    );
  }
});

module.exports = Output;
