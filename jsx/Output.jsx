/* jshint -W097,-W117 */
'use strict';

var React = require('react');
var _ = require('underscore');

var Table = require('./lib/reactable/table.jsx').Table;
var Thead = require('./lib/reactable/thead.jsx').Thead;
var Tfoot = require('./lib/reactable/tfoot.jsx').Tfoot;
var Th = require('./lib/reactable/th.jsx').Th;
var Tr = require('./lib/reactable/tr.jsx').Tr;
var Td = require('./lib/reactable/td.jsx').Td;

var ConnectionWindow = require('./ConnectionWindow.jsx');

class Output extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleCellClick = this.handleCellClick.bind(this);
    this.handleOnCellEdit = this.handleOnCellEdit.bind(this);


    this.state = {
      clickedCellPosition: {rowId: -1, colId: -1},
      editedCells: {},

      filterString: '',
    };
  }

  generateEditedCellKey(rowId, colId) {
    return rowId + "." + colId;
  }

  resultsToDictionary(results) {
    var reformatData = _.map(results.rows, function(row){
      var some = {};
      _.each(results.columns,function(key,i){some[key] = row[i];});
      return some;
    });

    return reformatData;
  }

  handleCellClick(rowId, colId, e) {
    console.log("handleCellClick ", rowId, colId);

    this.setState({
      clickedCellPosition: {rowId: rowId, colId: colId},
    });
  }

  handleOnCellEdit(rowId, colId, e) {
    console.log("handleOnCellEdit ", rowId, colId, e.target.value);

    var tempEditedCells = _.clone(this.state.editedCells);
    tempEditedCells[this.generateEditedCellKey(rowId, colId)] = e.target.value;

    console.log("Or: ", this.state.editedCells, "Changed", tempEditedCells);

    this.setState({
      editedCells: tempEditedCells,
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
        <Th key={i} data={col} column={col}>{col}</Th>
      );
    });

    return (
      <Thead>{children}</Thead>
    );
  }

  renderRow(row, rowId) {
    var self = this;
    var colId = -1;
    var children = _.map(row, function(value, col) {
      colId = colId + 1;

      var position = {rowId: rowId, colId: colId}

      if (self.state.clickedCellPosition.rowId == rowId && self.state.clickedCellPosition.colId == colId) {
        var isEditable = true
      }

      if (self.state.editedCells[self.generateEditedCellKey(rowId, colId)] != undefined) {
        var value = self.state.editedCells[self.generateEditedCellKey(rowId, colId)];
      }

      return (
        <Td
          key={position}
          column={col}
          position={position}
          onClick={self.handleCellClick.bind(self, rowId, colId)}
          isEditable={isEditable}
          onEdit={self.handleOnCellEdit.bind(self, rowId, colId)}>
            {value}
        </Td>
      );
    });

    return (
      <Tr key={rowId} >{children}</Tr>
    );
  }

  renderFooter() {
    return (
      <Tfoot>
        <tr className="foot">
          <td className="foot" colspan="99999">Temp Footer</td>
        </tr>
      </Tfoot>
    );
  }

  renderResults(results) {
    var data = this.resultsToDictionary(results);
    var header = this.renderHeader(results.columns);

    var self = this;
    var rows = _.map(data, function(row, i) {
      return self.renderRow(row, i);
    });

    var footer = this.renderFooter();

    if (this.props.withInput) {
      var filterable = results.columns;
      var filterPlaceholder = "Filter Results";
      var itemsPerPage = 100;
      var filterStyle = { top: this.props.dragBarPosition + 6 + 'px' };
    } else {
      var tableStyle = { height: '0' };
    }

    if (this.props.isSidebar) {
      return (
        <Table
          id="sidebar-modal-results"
          className="sidebar-modal-results"
          sortable={true} >
            {header}
            {rows}
        </Table>
      );
    }

    return (
      <Table
        id="results"
        className="results"
        style={tableStyle}
        sortable={true}
        filterable={filterable}
        filterPlaceholder={filterPlaceholder}
        filterStyle={filterStyle}
        onFilter={filter => {
            this.setState({ filterString: filter });
        }}
        filterString={this.state.filterString}
        itemsPerPage={itemsPerPage}
        resetPagination={this.props.resetPagination} >
          {header}
          {rows}
      </Table>
    );
  }

  renderNoResults() {
    return (
      <div>
          No records found
      </div>
    );
  }

  renderError(errorMsg) {
    return (
      <div>
          Err: {errorMsg}
      </div>
    );
  }

  render() {
    var clsOutput, children;
    var results = this.props.results;
    if (!results) {
      children = this.renderNoResults();
      clsOutput = "empty";
    } else {
      if (results.error) {
        children = this.renderError(results.error);
      } else if (!results.rows || results.rows.length === 0) {
        children = this.renderNoResults();
        clsOutput = "empty";
      } else {
        children = this.renderResults(results);
      }
    }

    if (this.props.isSidebar) {
      return (
        <div id="sidebar-result-wrapper">
          {children}
        </div>
      );
    }

    var outputStyle = { top: this.props.dragBarPosition + 60 + 'px'};
    if (clsOutput != "empty") {
      outputStyle['marginTop'] = '-10px';
    }

    if (!this.props.withInput) {
      outputStyle['top'] = '60px';
    }

    return (
      <div id="output" className={clsOutput} style={outputStyle}>
        <div id="wrapper">
          {children}
        </div>
      </div>
    );
  }
}

module.exports = Output;