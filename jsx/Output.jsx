import React from 'react';
import { Table } from './reactable/table.jsx';
import { Thead } from './reactable/thead.jsx';
import { Tfoot } from './reactable/tfoot.jsx';
import { Th } from './reactable/th.jsx';
import { Tr } from './reactable/tr.jsx';
import { Td } from './reactable/td.jsx';
import ConnectionWindow from './ConnectionWindow.jsx';
import QueryEditBar from './QueryEditBar.jsx';
import * as action from './action.js';
import view from './view.js';

function resultsToDictionary(results) {
  const reformatData = results.rows.map(function(row) {
    let some = {};
    results.columns.forEach((key,i) => some[key] = row[i]);
    return some;
  });

  return reformatData;
}

export default class Output extends React.Component {
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
    var tmp = Object.assign({}, this.props.editedCells);
    if (tmp[rowId] == undefined) {
      tmp[rowId] = {};
    }
    tmp[rowId][colId] = value;

    action.editedCells(tmp);
  }

  getEditedCells(rowId, colId) {
    if (this.props.editedCells[rowId] == undefined) {
      return undefined;
    }
    return this.props.editedCells[rowId][colId];
  }

  generateQuery() {
    const table = this.props.selectedTable;
    const results = this.props.results;
    const tableStructures = this.props.tableStructures;
    const resultsAsDictionary = resultsToDictionary(results);
    const editedCells = this.props.editedCells;

    let query = "";
    for (let rowId in editedCells) {
      let value = editedCells[rowId];
      //let values = rowId.split('.');

      const thisRow = editedCells[rowId];
      let index = 0;
      let colsAfterEdit = "";
      for (let colId in thisRow) {
        let value = thisRow[colId];
        const columnToBeEdited = results.columns[colId];
        const afterChange = value;

        if (afterChange == "") {
          colsAfterEdit += columnToBeEdited + "=NULL ";
        } else {
          colsAfterEdit += columnToBeEdited + "=\'" + afterChange + "\'";
        }

        if (index < Object.keys(thisRow).length - 1) {
          colsAfterEdit += ", ";
        }
        index += 1;
      }

      const columns = results.columns.join(", ");

      const tableStructuresAsDictionary = resultsToDictionary(tableStructures[table]);
      let schema = null;
      if (tableStructuresAsDictionary.length > 0) {
        schema = tableStructuresAsDictionary[0]["table_schema"];
      } else {
        console.log("THIS CASE SHOULD NOT HAPPEN IS THERE A WAY TO LOG THIS?");
      }

      const rowAsDictionary = resultsAsDictionary[rowId];

      index = 0;
      let rowToBeEdited = "";
      for (let key in rowAsDictionary) {
        value = rowAsDictionary[key];
        if (value == null) {
          rowToBeEdited += key + " IS NULL ";
        } else {
          rowToBeEdited += key + "=\'" + value + "\' ";
        }

        if (index < Object.keys(rowAsDictionary).length - 1) {
          rowToBeEdited += "AND ";
        }
        index += 1;
      }

      query += `UPDATE ${schema}.${table}
SET ${colsAfterEdit}
WHERE ctid IN (SELECT ctid FROM ${schema}.${table}
WHERE ${rowToBeEdited}
LIMIT 1 FOR UPDATE)
RETURNING ${columns};
`;

      console.log("QUERY:", query);
      // WHERE countrycode='ABW' AND language='Not English no qq' AND isofficial='false' AND percentage='9.5'
      // UPDATE countrylanguage SET language='Not furkan' WHERE ctid IN (SELECT ctid FROM countrylanguage WHERE countrycode='ABW' AND language='Not English no qq' AND isofficial='false' AND percentage='9.5' LIMIT 1 FOR UPDATE) RETURNING language;

    }

    return query;
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
    let i = 0;
    columns = columns || [];
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
    let colId = -1;
    const selectedCellPosition = this.props.selectedCellPosition;
    const selectedView = this.props.selectedView;
    let children = [];
    for (let col in row) {
      let value = row[col];
      colId = colId + 1;
      const position = {rowId: rowId, colId: colId};

      const isEditable = selectedCellPosition.rowId == rowId &&
          selectedCellPosition.colId == colId &&
          selectedView == view.SQLQuery;

      let tdStyle = {};
      if (this.getEditedCells(rowId, colId) != undefined) {
        value = this.getEditedCells(rowId, colId);
        tdStyle = {
          background: '#7DCED2',
          color: '#ffffff',
          // border: 'solid 1px #3B8686',
        };
      }

      children.push(
        <Td
          key={position}
          column={col}
          position={position}
          style={tdStyle}
          onClick={this.handleCellClick.bind(this, rowId, colId)}
          isEditable={isEditable}
          onEdit={this.handleOnCellEdit.bind(this, rowId, colId)}>
            {value}
        </Td>
      );
    }

    return (
      <Tr key={rowId}>{children}</Tr>
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
    const data = resultsToDictionary(results);
    const header = this.renderHeader(results.columns);
    const rows = data.map((row, i) => this.renderRow(row, i));
    //const footer = this.renderFooter();

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
