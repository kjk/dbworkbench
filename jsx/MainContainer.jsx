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
    this.onMouseDown = this.onMouseDown.bind(this);
    this.onMouseMove = this.onMouseMove.bind(this);
    this.onMouseUp = this.onMouseUp.bind(this);

    this.sidebarDx = store.getSidebarDx();
    this.state = {
      dragBarPosition: 200,
      dragging: false,
    };
  }

  componentWillMount() {
    this.cidSidebarDx = store.onSidebarDx( (dx) => {
      this.sidebarDx = dx;
      const el = ReactDOM.findDOMNode(this);
      el.style.left = dx + "px";
    });
  }

  componentWillUnmount() {
    store.offSidebarDx(this.cidSidebarDx);
  }

  componentDidUpdate(props, state) {
    if (this.state.dragging && !state.dragging) {
      document.addEventListener('mousemove', this.onMouseMove);
      document.addEventListener('mouseup', this.onMouseUp);
    } else if (!this.state.dragging && state.dragging) {
      document.removeEventListener('mousemove', this.onMouseMove);
      document.removeEventListener('mouseup', this.onMouseUp);
    }
  }

  onMouseDown(e) {
    // only left mouse button
    if (e.button !== 0) return;
    this.setState({
      dragging: true,
    });
    e.stopPropagation();
    e.preventDefault();
  }

  onMouseUp(e) {
    this.setState({
      dragging: false,
    });
    e.stopPropagation();
    e.preventDefault();
  }

  onMouseMove(e) {
    const minDragbarDx = 60;
    const maxDragbarDx = 400;

    if (!this.state.dragging) return;
    if ((e.pageY < minDragbarDx) || (e.pageY > maxDragbarDx)) {
      return;
    }
    this.setState({
      dragBarPosition: e.pageY,
    });
    e.stopPropagation();
    e.preventDefault();
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
      <div id="body" style={style}>
          <DbNav view={this.props.selectedView}/>

          { withInput ?
            <Input
              dragBarPosition={this.state.dragBarPosition}
              supportsExplain={this.props.supportsExplain}
              onMouseDown={this.onMouseDown}
              onMouseMove={this.onMouseMove}
              onMouseUp={this.onMouseUp}
              editedCells={this.props.editedCells} />
            : null
          }
          <Output
            dragBarPosition={this.state.dragBarPosition}
            selectedView={this.props.selectedView}
            results={this.props.results}
            withInput={withInput}
            resetPagination={this.props.resetPagination}
            tableStructures={this.props.tableStructures}
            selectedTable={this.props.selectedTable}
            selectedCellPosition={this.props.selectedCellPosition}
            editedCells={this.props.editedCells} />
      </div>
    );
  }
}
