var React = require('react');
var DbNav = require('./DbNav.jsx');
var Input = require('./Input.jsx');
var Output = require('./Output.jsx');
var view = require('./view.js');

class MainContainer extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.onMouseDown = this.onMouseDown.bind(this);
    this.onMouseMove = this.onMouseMove.bind(this);
    this.onMouseUp = this.onMouseUp.bind(this);

    this.state = {
      dragBarPosition: 200,
      dragging: false,
    };
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
    // when showing sql query, results are below editor window
    var withInput = (this.props.selectedView === view.SQLQuery);

    var divStyle = {
      left: this.props.dragBarPosition + 'px',
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
      <div id="body" style={divStyle}>
          <DbNav view={this.props.selectedView}/>

          { withInput ?
            <Input
              dragBarPosition={this.state.dragBarPosition}
              supportsExplain={this.props.supportsExplain}
              onMouseDown={this.onMouseDown}
              onMouseMove={this.onMouseMove}
              onMouseUp={this.onMouseUp}
              spinnerVisible={this.props.spinnerVisible} />
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

module.exports = MainContainer;