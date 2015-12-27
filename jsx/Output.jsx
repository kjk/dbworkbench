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
var QueryEditBar = require('./QueryEditBar.jsx');

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

  componentWillReceiveProps(nextProps) {
    if (nextProps.resetPagination) { // TODO: Maybe use another name instead of resetPagination
      this.setState({
        clickedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
        filterString: '',
      });
    }
  }

  generateEditedCellKey(rowId, colId) {
    return rowId + "." + colId;
  }

  generateQuery() {
    var self = this;
    var table = this.props.selectedTable;
    var query = "";
    var resultsAsDictionary = this.resultsToDictionary(this.props.results);

    _.each(this.state.editedCells, function(value, key, obj) {
      var values = key.split('.');
      var rowId = values[0];
      var colId = values[1];

      var columnToBeEdited = self.props.results.columns[colId];
      var afterChange = value;
      var columns = self.props.results.columns.join(", ");

      var tableStructuresAsDictionary = self.resultsToDictionary(self.props.tableStructures[table]);
      if (tableStructuresAsDictionary.length > 0) {
        var schema = tableStructuresAsDictionary[0]["table_schema"];
      } else {
        console.log("THIS CASE SHOULD NOT HAPPEN IS THERE A WAY TO LOG THIS?");
      }

      var rowAsDictionary = resultsAsDictionary[rowId];

      var index = 0;
      var rowToBeEdited = "";
      _.each(rowAsDictionary, function(value, key, obj) {
        rowToBeEdited += key + "=\'" + obj[key] + "\' ";
        if (index < Object.keys(rowAsDictionary).length - 1) {
          rowToBeEdited += "AND ";
        }
        index += 1;
      });

      query += "UPDATE " + schema + "." + table + " ";
      query += "SET " + columnToBeEdited + "=\'" + afterChange + "\' ";
      query += "WHERE ctid IN (SELECT ctid FROM " + schema + "." + table + " ";
      query += "WHERE " + rowToBeEdited + " ";
      query += "LIMIT 1 FOR UPDATE) ";
      query += "RETURNING " + columns + ";";

      console.log("QUERY:", query);
      // WHERE countrycode='ABW' AND language='Not English no qq' AND isofficial='false' AND percentage='9.5'
      // UPDATE countrylanguage SET language='Not furkan' WHERE ctid IN (SELECT ctid FROM countrylanguage WHERE countrycode='ABW' AND language='Not English no qq' AND isofficial='false' AND percentage='9.5' LIMIT 1 FOR UPDATE) RETURNING language;

    });

    return query;
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

  handleDiscardChanges() {
    this.setState({
      clickedCellPosition: {rowId: -1, colId: -1},
      editedCells: {},
    });
  }

  handleOnCellEdit(rowId, colId, e) {
    console.log("handleOnCellEdit ", rowId, colId, e.target.value);

    var tempEditedCells = _.clone(this.state.editedCells);
    tempEditedCells[this.generateEditedCellKey(rowId, colId)] = e.target.value;

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

      var position = {rowId: rowId, colId: colId};

      if (self.state.clickedCellPosition.rowId == rowId && self.state.clickedCellPosition.colId == colId) {
        var isEditable = true;
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
        resetPagination={this.props.viewChanged} >
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

    var numberOfRowsEdited = Object.keys(this.state.editedCells).length;
    if (numberOfRowsEdited !== 0) {
      var queryEditBar = (
        <QueryEditBar
          numberOfRowsEdited={numberOfRowsEdited}
          generateQuery={this.generateQuery.bind(this)}
          onHandleDiscardChanges={this.handleDiscardChanges.bind(this)} />
        );
    }

    return (
      <div id="output" className={clsOutput} style={outputStyle}>
        <div id="wrapper">
          {children}
          {queryEditBar}
        </div>
      </div>
    );
  }
}

module.exports = Output;