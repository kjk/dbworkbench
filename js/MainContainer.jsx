import React from 'react';
import ReactDOM from 'react-dom';
import DbNav from './DbNav.jsx';
import DragBarHoriz from './DragBarHoriz.jsx';
import Input from './Input.jsx';
import Output from './Output2.jsx';
import view from './view.js';
import * as store from './store.js';

export default class MainContainer extends React.Component {
  constructor(props, context) {
    super(props, context);
  }

  componentWillMount() {
    store.onSidebarDx((dx) => {
      const el = ReactDOM.findDOMNode(this);
      el.style.left = dx + 'px';
    }, this);
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  handlePosChanged(y) {
    store.setQueryEditDy(y);
  }

  render() {
    // TODO: after database connect, this happens 28 times
    //console.log("MainContainer render");

    const withInput = (this.props.selectedView === view.SQLQuery);

    const style = {
      left: store.getSidebarDx()
    };

    return (
      <div id="body" style={ style }>
        <DbNav view={ this.props.selectedView } />
        { withInput ?
          <Input supportsExplain={ this.props.supportsExplain } editedCells={ this.props.editedCells } />
          : null }
        { withInput ?
          <DragBarHoriz initialY={ store.getQueryEditDy() }
            min={ 60 }
            max={ 400 }
            onPosChanged={ this.handlePosChanged } />
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

MainContainer.propTypes = {
  selectedView: React.PropTypes.string,
  supportsExplain: React.PropTypes.bool,
  results: React.PropTypes.any, // TODO: more specific,
  resetPagination: React.PropTypes.bool, // TODO: more specific
  tableStructures: React.PropTypes.any, // TODO: more specific
  selectedTable: React.PropTypes.string,
  selectedCellPosition: React.PropTypes.any, // TODO: more specific
  editedCells: React.PropTypes.any // TODO: more specific
};
