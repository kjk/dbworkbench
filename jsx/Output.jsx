import React from 'react';
import ReactDOM from 'react-dom';
import { Table } from './reactable/table.jsx';
import { Thead } from './reactable/thead.jsx';
import { Tr } from './reactable/tr.jsx';
import { Td } from './reactable/td.jsx';
import QueryEditBar from './QueryEditBar.jsx';
import view from './view.js';
import * as action from './action.js';
import * as store from './store.js';

function resultsToDictionary(results) {
  const reformatData = results.rows.map(function(row) {
    let some = {};
    results.columns.forEach((key, i) => some[key] = row[i]);
    return some;
  });

  return reformatData;
}

export default class Output extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleCellClick = this.handleCellClick.bind(this);
    this.handleOnCellEdit = this.handleOnCellEdit.bind(this);

    this.queryEditDy = store.getQueryEditDy();
  }

  componentWillMount() {
    store.onQueryEditDy((dy) => {
      this.queryEditDy = dy;
      const el = ReactDOM.findDOMNode(this);
      el.style.top = this.topPos();
    }, this);
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.resetPagination) { // TODO: Maybe use another name instead of resetPagination
      action.clearFilter();
    }
  }

  setEditedCells(rowId, colId, value) {
    let tmp = Object.assign({}, this.props.editedCells);
    if (tmp[rowId] == undefined) {
      tmp[rowId] = {};
    }
    tmp[rowId][colId] = value;

    action.editedCells(tmp);
  }

  getEditedCells(rowId, colId) {
    if (this.props.editedCells == undefined || this.props.editedCells[rowId] == undefined) {
      return undefined;
    }
    return this.props.editedCells[rowId][colId];
  }

  generateQuery() {
    const tableName = this.props.selectedTable;
    const results = this.props.results;
    const tableStructures = this.props.tableStructures;
    const resultsAsDictionary = resultsToDictionary(results);
    const editedCells = this.props.editedCells;
    const schema = this.props.tableStructures[tableName]['table_schema'];

    let query = '';
    for (let rowId in editedCells) {
      let value = editedCells[rowId];
      //let values = rowId.split('.');

      const thisRow = editedCells[rowId];
      let index = 0;
      let colsAfterEdit = '';
      for (let colId in thisRow) {
        let value = thisRow[colId];
        const columnToBeEdited = results.columns[colId];
        const afterChange = value;

        if (afterChange == '') {
          colsAfterEdit += columnToBeEdited + '=NULL ';
        } else {
          colsAfterEdit += columnToBeEdited + "='" + afterChange + "'";
        }

        if (index < Object.keys(thisRow).length - 1) {
          colsAfterEdit += ', ';
        }
        index += 1;
      }

      const columns = results.columns.join(', ');
      const rowAsDictionary = resultsAsDictionary[rowId];

      index = 0;
      let rowToBeEdited = '';
      for (let key in rowAsDictionary) {
        value = rowAsDictionary[key];
        if (value == null) {
          rowToBeEdited += key + ' IS NULL ';
        } else {
          rowToBeEdited += key + "='" + value + "' ";
        }

        if (index < Object.keys(rowAsDictionary).length - 1) {
          rowToBeEdited += 'AND ';
        }
        index += 1;
      }

      query += `UPDATE ${schema}.${tableName}
SET ${colsAfterEdit}
WHERE ctid IN (SELECT ctid FROM ${schema}.${tableName}
WHERE ${rowToBeEdited}
LIMIT 1 FOR UPDATE)
RETURNING ${columns};
`;

      console.log('QUERY:', query);
      // WHERE countrycode='ABW' AND language='Not English no qq' AND isofficial='false' AND percentage='9.5'
      // UPDATE countrylanguage SET language='Not furkan' WHERE ctid IN (SELECT ctid FROM countrylanguage WHERE countrycode='ABW' AND language='Not English no qq' AND isofficial='false' AND percentage='9.5' LIMIT 1 FOR UPDATE) RETURNING language;

    }

    return query;
  }

  handleCellClick(e) {
    const node = e.target.nodeName;
    if (node != 'TD') {
      return;
    }
    const rowColStr = e.target.attributes['data-custom-attribute'].value;
    const parts = rowColStr.split('-');
    if (parts.length != 2) {
      return;
    }
    const rowId = parseInt(parts[0], 10);
    const colId = parseInt(parts[1], 10);
    console.log('handleCellClick rowId: ', rowId, ' : ', colId);
    action.selectedCellPosition({
      rowId: rowId,
      colId: colId
    });
  }

  handleDiscardChanges() {
    // TODO: do these togethor
    action.editedCells({});
    action.selectedCellPosition({
      rowId: -1,
      colId: -1
    });
  }

  handleOnCellEdit(rowId, colId, e) {
    console.log('handleOnCellEdit ', rowId, colId, e.target.value);
    this.setEditedCells(rowId, colId, e.target.value);
  }

  renderHeader(columns, sortColumn, sortOrder) {
    let i = 0;
    columns = columns || [];
    var children = columns.map(function(col) {
      // TODO: use sortColumn and sortOrder)
      i = i + 1;
      return (
        <th key={ i } data={ col } column={ col }>
          { col }
        </th>
        );
    });

    return (
      <Thead>
        { children }
      </Thead>
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
      const position = {
        rowId: rowId,
        colId: colId
      };

      var isEditable = false;
      if (selectedCellPosition != undefined) {
        isEditable = selectedCellPosition.rowId == rowId &&
          selectedCellPosition.colId == colId &&
          selectedView == view.SQLQuery;
      }

      let tdStyle = {};
      if (this.getEditedCells(rowId, colId) != undefined) {
        value = this.getEditedCells(rowId, colId);
        tdStyle = {
          background: '#6EACE3',
          color: '#ffffff',
        // border: 'solid 1px #3B8686',
        };
      }

      children.push(
        <Td key={ position }
          column={ col }
          position={ position }
          style={ tdStyle }
          isEditable={ isEditable }
          onEdit={ this.handleOnCellEdit.bind(this, rowId, colId) }>
          { value }
        </Td>
      );
    }

    return (
      <Tr key={ rowId }>
        { children }
      </Tr>
      );
  }

  renderResults(results) {
    const data = resultsToDictionary(results);
    const header = this.renderHeader(results.columns);
    const rows = data.map((row, i) => this.renderRow(row, i));

    let nEdited = 0;
    if (this.props.editedCells) {
      nEdited = Object.keys(this.props.editedCells).length;
    }
    let filterable = [];
    if (this.props.withInput && nEdited == 0) {
      filterable = results.columns;
    }

    if (this.props.withInput) {
      var itemsPerPage = 100;
    } else {
      var tableStyle = {
        height: '0'
      };
    }

    if (this.props.isSidebar) {
      return (
        <Table id="sidebar-modal-results"
          className="sidebar-modal-results"
          onClick={ this.handleCellClick }
          sortable>
          { header }
          { rows }
        </Table>
        );
    }

    return (
      <Table id="results"
        className="results"
        onClick={ this.handleCellClick }
        style={ tableStyle }
        sortable
        filterable={ filterable }
        itemsPerPage={ itemsPerPage }
        resetPagination={ this.props.resetPagination }>
        { header }
        { rows }
      </Table>
      );
  }

  topPos() {
    let top = 60;
    if (this.props.withInput) {
      top = this.queryEditDy + 60;
    }
    return top + 'px';
  }

  renderEmptyOrError(results) {
    let res;
    if (results && results.error) {
      res = <div>
              Error:
              { results.error }
            </div>;
    } else if (!results || !results.rows || results.rows.length == 0) {
      res = <div>
              No records found
            </div>;
    }
    if (!res) {
      return res;
    }
    if (this.props.isSidebar) {
      return (<div id="sidebar-result-wrapper">
                { res }
              </div>);
    }
    let style = {
      top: this.topPos()
    };

    return (<div id="output" className="empty" style={ style }>
              <div id="wrapper">
                { res }
              </div>
            </div>
      );
  }

  render() {
    //console.log("Output.render");

    const results = this.props.results;
    const res = this.renderEmptyOrError(results);
    if (res) {
      return res;
    }

    const children = this.renderResults(results);

    if (this.props.isSidebar) {
      return (
        <div id="sidebar-result-wrapper">
          { children }
        </div>
        );
    }

    const editedCells = this.props.editedCells || {};
    const nEdited = Object.keys(editedCells).length;
    const showQueryBar = nEdited > 0;

    let style = {
      top: this.topPos(),
      marginTop: -10
    };

    return (
      <div id="output" style={ style }>
        <div id="wrapper">
          { children }
          { showQueryBar ?
            <QueryEditBar numberOfRowsEdited={ nEdited } generateQuery={ this.generateQuery.bind(this) } onHandleDiscardChanges={ this.handleDiscardChanges.bind(this) } />
            : null }
        </div>
      </div>
      );
  }
}
