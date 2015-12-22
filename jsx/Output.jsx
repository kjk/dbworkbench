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
    this.handleRowClick = this.handleRowClick.bind(this);

    this.state = {
      clickedRowKey: -1,
      rowStyle: {},

      filterString: '',

    };
  }

  resultsToDictionary(results) {
    var reformatData = _.map(results.rows, function(row){
      var some = {};
      _.each(results.columns,function(key,i){some[key] = row[i];});
      return some;
    });

    return reformatData;
  }

  handleRowClick(key, e) {
    console.log("Enlarging ", key);
    var enlargeStyle = {
      maxWidth: '350px',
      maxHeight: '100%',
      overflow: 'hidden',
      textOverflow: 'ellipsis',
      whiteSpace: 'normal',
    };

    if (_.isEqual(enlargeStyle, this.state.rowStyle) && key == this.state.clickedRowKey) {
      console.log("Shrinking");
      enlargeStyle = {};
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
        <Th key={i} data={col} column={col}>{col}</Th>
      );
    });

    return (
      <Thead>{children}</Thead>
    );
  }

  renderRow(row, key) {
    var style = {};
    if (this.state.clickedRowKey == key) {
      style = this.state.rowStyle;
    }

    var i = 0;
    var children = _.map(row, function(row, col) {
      i = i + 1;

      // console.log("row", row, "col", col)
      return (
        <Td key={i} column={col} value={row}><div style={style}>{row}</div></Td>
      );
    });

    return (
      <Tr key={key} onClick={this.handleRowClick.bind(this,key)}>{children}</Tr>
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