import React from 'react';
import ReactDOM from 'react-dom';
import ResultsPaginator from './ResultsPaginator.jsx';
import QueryEditBar from './QueryEditBar.jsx';
import view from './view.js';
import * as action from './action.js';
import * as store from './store.js';
import * as sort from './reactable/sort.jsx';

/*
Output selectedView="Query", results={ columns: ["colName1", "colName2", ...]
       rows: [ [], [] ]} resetPagination=true|false
       selectedTable="address" withInput=true|false
       tableStructures=[]
       selectedCellPosition={colId: 1, rowId: 6}
       editedCells: { 6: { 1: "new value" }}
  div id="output" style={top: , marginTop}
    div id="wrapper"
      Table children=[Thead, [[Tr], [Tr], ...]] className="results" defaultSort=false
            filterBy="" id="results" itemsPerPage=100 onClick=func(), resetPagination=true
            sortBy=false, sortable=true
            filterable=["colName1", "colName2"]
        div
          table children, className, filterBy, id, onClick, resetPagination
            Thead columns=[ {key: "colName", label: "colName"}, ...], onSort
                  sort={column: null, direction: 1} sortableColumns={ colName1: "default", colName2: "default", ...}
              tr children, className="reactable-column-header"
                th children="colName"
                   className="reactable-th-$colName reactable-header-sortable"
                   onClick, role="button", tabIndex=0
                ...
            tbody children className="reactable-data"
              Tr columns= data=
                tr children
                  Td children=1, column={key: label:}, isEditable=false onEdit
                     position={colId: rowId: } style={}
                    td children="1" label=$colName data-custom-attribute="0-0"
                  ..
      Paginator
*/
class ColumnInfo {
  constructor(name, sortOrder) {
    this.name = name;
    this.sortOrder = sortOrder;
    this.isSortable = true;
  }
}

function incSortOrder(sortOrder) {
  if (sortOrder == sort.None) {
    return sort.Up;
  }
  if (sortOrder == sort.Up) {
    return sort.Down;
  }
  return sort.Up;
}

function calcColumnInfos(columnNames, sortByColumnIdx, prevColumnInfos) {
  if (!columnNames) {
    return [];
  }
  const res = columnNames.map((name, idx) => {
    let sortOrder = sort.None;
    if (prevColumnInfos && idx == sortByColumnIdx) {
      const prevSortOrder = prevColumnInfos[sortByColumnIdx].sortOrder;
      sortOrder = incSortOrder(prevSortOrder);
    }
    return new ColumnInfo(name, sortOrder);
  });
  return res;
}

function topPos(dy, withInput) {
  let top = withInput ? dy + 60 : 60;
  return top + 'px';
}

const nPerPage = 100;

function getPage(arr, pageNo) {
  const start = pageNo * nPerPage;
  let end = start + nPerPage;
  if (end > arr.length) {
    end = arr.length;
  }
  let res = [];
  for (let i = start; i < end; i++) {
    res.push(arr[i]);
  }
  return res;
}

export default class Output extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.handleDiscardChanges = this.handleDiscardChanges.bind(this);
    this.handleColumnClick = this.handleColumnClick.bind(this);

    this.state = this.calcState(this.props);
  }

  componentWillMount() {
    //console.log('Output2.componentWillMount');
    store.onQueryEditDy(dy => {
      const el = ReactDOM.findDOMNode(this);
      el.style.top = topPos(dy, this.props.withInput);
    }, this);
  }

  componentWillReceiveProps(nextProps) {
    //console.log('Output2.componentWillReceiveProps');
    this.setState(this.calcState(nextProps));
  }

  componentWillUnmount() {
    //console.log('Output2.componentWillUnmount');
    store.offAllForOwner(this);
  }

  handlePageChanged(pageNo) {
    //console.log('Output2.handlePageChanged: ', pageNo);
    const rows = getPage(this.state.allRows, pageNo);
    this.setState({
      currPageNo: pageNo,
      rows: rows
    });
  }

  handleDiscardChanges() {
    //console.log('Output2.handleDiscardChanges');
    // TODO: do these togethor
    /*
    action.editedCells({});
    action.selectedCellPosition({
      rowId: -1,
      colId: -1
    });*/
  }

  handleColumnClick(e) {
    e.preventDefault();
    const colIdx = e.target.cellIndex;
    console.log('Output2.handleColumnClick: ', colIdx);
    const columns = this.state.columns;
    const columnInfos = calcColumnInfos(columns, colIdx, this.state.columnInfos);
    this.setState({
      columnInfos: columnInfos
    });
  }

  calcState(props) {
    this.top = topPos(store.getQueryEditDy(), this.props.withInput);

    //console.log('Output2.calcState');
    const results = props.results;
    const allRows = results ? results.rows : [];
    const columns = results ? results.columns : [];
    const nPages = Math.ceil(allRows.length / nPerPage);
    const pageNo = 0;
    const rows = getPage(allRows, pageNo);
    let columnInfos = this.state ? this.state.columnInfos : null;
    if (!columnInfos || columnInfos.length != columns.length) {
      columnInfos = calcColumnInfos(columns, 0, null);
    }
    return {
      results: results,
      allRows: allRows,
      currPageNo: pageNo,
      nPages: nPages,
      rows: rows,
      columns: columns,
      columnInfos: columnInfos
    };
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
      top: this.top
    };

    return (
      <div id="output" className="empty" style={ style }>
        <div id="wrapper">
          { res }
        </div>
      </div>
      );
  }

  renderTd(rowIdx, colIdx, colData) {
    const key = '' + rowIdx + '-' + colIdx;
    return (
      <td key={ key } data-custom-attribute={ key }>
        { colData }
      </td>
      );
  }

  renderTr(rowIdx, row) {
    const children = row.map((col, colIdx) => this.renderTd(rowIdx, colIdx, col));
    return (
      <tr key={ rowIdx }>
        { children }
      </tr>
      );
  }

  renderTheadTh(col, colIdx) {
    let cls = 'reactable-header-sortable';
    const s = col.name;
    if (col.sortOrder == sort.Up) {
      cls += ' reactable-header-sort-asc';
    } else if (col.sortOrder == sort.Down) {
      cls += ' reactable-header-sort-desc';
    }
    return (
      <th key={ colIdx }
        className={ cls }
        role="button"
        tabIndex="0"
        onClick={ this.handleColumnClick }>
        { s }
      </th>
      );
  }

  renderResults() {
    const state = this.state;
    const results = state.results;
    const allRows = state.allRows;
    const columns = state.columnInfos.map((col, colIdx) => this.renderTheadTh(col, colIdx));
    const rows = state.rows.map((row, rowIdx) => this.renderTr(rowIdx, row));
    const pageChanged = pageNo => this.handlePageChanged(pageNo);
    return (
      <div>
        <table className="results" id="results">
          <thead>
            <tr className="reactable-column-header">
              { columns }
            </tr>
          </thead>
          <tbody className="reactable-data">
            { rows }
          </tbody>
        </table>
        <ResultsPaginator nRows={ allRows.length }
          nPages={ this.state.nPages }
          currentPage={ this.state.currPageNo }
          onPageChange={ pageChanged } />
      </div>
      );
  }

  render() {
    console.log('Output2.render');

    const res = this.renderEmptyOrError(this.state.results);
    if (res) {
      return res;
    }

    const children = this.renderResults();

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
      top: this.top,
      marginTop: -10
    };

    return (
      <div id="output" style={ style }>
        <div id="wrapper">
          { children }
          { showQueryBar ?
            <QueryEditBar numberOfRowsEdited={ nEdited } onHandleDiscardChanges={ this.handleDiscardChanges } />
            : null }
        </div>
      </div>
      );
  }
}

Output.propTypes = {
  withInput: React.PropTypes.bool,
  isSidebar: React.PropTypes.bool,
  editedCells: React.PropTypes.any // TODO: be more specific
};
