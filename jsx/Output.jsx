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
var action = require('./action.js');
var view = require('./view.js');

class Output extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleCellClick = this.handleCellClick.bind(this);
    this.handleOnCellEdit = this.handleOnCellEdit.bind(this);

    this.state = {
      filterString: '',
    };
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.resetPagination) { // TODO: Maybe use another name instead of resetPagination
      this.setState({
        filterString: '',
      });
    }
  }

  setEditedCells(rowId, colId, value) {
    var tempEditedCells = _.clone(this.props.editedCells);
    if (tempEditedCells[rowId] == undefined) {
      tempEditedCells[rowId] = {};
    }
    tempEditedCells[rowId][colId] = value;

    action.editedCells(tempEditedCells);
  }

  getEditedCells(rowId, colId) {
    if (this.props.editedCells[rowId] == undefined) {
      return undefined;
    }
    return this.props.editedCells[rowId][colId];
  }

  generateQuery() {
    var self = this;
    var table = this.props.selectedTable;
    var query = "";
    var resultsAsDictionary = this.resultsToDictionary(this.props.results);

    _.each(this.props.editedCells, function(value, key, obj) {
      var values = key.split('.');
      var rowId = key;

      var thisRow = obj[key];
      var index = 0;
      var colsAfterEdit = "";
      _.each(thisRow, function(value, key, obj) {
        var colId = key;
        var columnToBeEdited = self.props.results.columns[colId];
        var afterChange = value;

        if (afterChange == "") {
          colsAfterEdit += columnToBeEdited + "=NULL ";
        } else {
          colsAfterEdit += columnToBeEdited + "=\'" + afterChange + "\'";
        }

        if (index < Object.keys(thisRow).length - 1) {
          colsAfterEdit += ", ";
        }
        index += 1;
      });

      var columns = self.props.results.columns.join(", ");

      var tableStructuresAsDictionary = self.resultsToDictionary(self.props.tableStructures[table]);
      if (tableStructuresAsDictionary.length > 0) {
        var schema = tableStructuresAsDictionary[0]["table_schema"];
      } else {
        console.log("THIS CASE SHOULD NOT HAPPEN IS THERE A WAY TO LOG THIS?");
      }

      var rowAsDictionary = resultsAsDictionary[rowId];

      index = 0;
      var rowToBeEdited = "";
      _.each(rowAsDictionary, function(value, key, obj) {
        if (value == null) {
          rowToBeEdited += key + " IS NULL ";
        } else {
          rowToBeEdited += key + "=\'" + value + "\' ";
        }

        if (index < Object.keys(rowAsDictionary).length - 1) {
          rowToBeEdited += "AND ";
        }
        index += 1;
      });

      query += `UPDATE ${schema}.${table}
SET ${colsAfterEdit}
WHERE ctid IN (SELECT ctid FROM ${schema}.${table}
WHERE ${rowToBeEdited}
LIMIT 1 FOR UPDATE)
RETURNING ${columns};`;

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
    action.selectedCellPosition({rowId: rowId, colId: colId});
  }

  handleDiscardChanges() {
    // TODO: do these togethor
    action.editedCells({});
    action.selectedCellPosition({rowId: -1, colId: -1});
  }

  handleOnCellEdit(rowId, colId, e) {
    console.log("handleOnCellEdit ", rowId, colId, e.target.value);
    this.setEditedCells(rowId, colId, e.target.value);
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

      if (self.props.selectedCellPosition.rowId == rowId &&
          self.props.selectedCellPosition.colId == colId &&
          self.props.selectedView == view.SQLQuery) {
        var isEditable = true;
      }

      if (self.getEditedCells(rowId, colId) != undefined) {
        var value = self.getEditedCells(rowId, colId);
        var tdStyle = {
          background: '#7DCED2',
          color: '#ffffff',
          border: 'solid 1px #3B8686',
        };
      }

      return (
        <Td
          key={position}
          column={col}
          position={position}
          style={tdStyle}
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

    var numberOfRowsEdited = Object.keys(this.props.editedCells).length;
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