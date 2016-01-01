import React from 'react';
import ReactDOM from 'react-dom';
import DbNav from './DbNav.jsx';
import Input from './Input.jsx';
import Output from './Output.jsx';
import view from './view.js';
import * as store from './store.js';

export default class MainContainer extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.sidebarDx = store.getSidebarDx();
  }

  componentWillMount() {
    store.onSidebarDx((dx) => {
      this.sidebarDx = dx;
      const el = ReactDOM.findDOMNode(this);
      el.style.left = dx + 'px';
    }, this);
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  // renderInput(tooLong, supportsExplain, inputStyle) {
  //   if (this.props.selectedView === view.SQLQuery) {
  //     return <Input style={inputStyle} tooLong={tooLong} supportsExplain={supportsExplain}/>;
  //   }
  // }

  render() {
    // TODO: after database connect, this happens 28 times
    //console.log("MainContainer render");

    // when showing sql query, results are below editor window
    var withInput = (this.props.selectedView === view.SQLQuery);

    var style = {
      left: this.sidebarDx,
    };

    // var results = this.props.results
    // if (results != null && results.rows != null) {
    //   if (results.rows.length > 100) {
    //     // It's only showed when +100. We could make this default.
    //     var tooLong = "Showing 100 out of " + results.rows.length + " rows."
    //     results.rows = results.rows.slice(0, 100);
    //   }
    // }

    return (
      <div id="body" style={ style }>
        <DbNav view={ this.props.selectedView } />
        { withInput ?
          <Input supportsExplain={ this.props.supportsExplain } editedCells={ this.props.editedCells } />
          : null }
        <Output selectedView={ this.props.selectedView }
          results={ this.props.results }
          withInput={ withInput }
          resetPagination={ this.props.resetPagination }
          tableStructures={ this.props.tableStructures }
          selectedTable={ this.props.selectedTable }
          selectedCellPosition={ this.props.selectedCellPosition }
          editedCells={ this.props.editedCells } />
      </div>
      );
  }
}
