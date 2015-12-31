import 'babel-polyfill';
import React from 'react';
import ReactDOM from 'react-dom';
import * as store from './store.js';

// TODO: add propTypes
// http://www.newmediacampaigns.com/blog/refactoring-react-components-to-es6-classes
export default class DragBarVert extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.onMouseDown = this.onMouseDown.bind(this);
    this.onMouseMove = this.onMouseMove.bind(this);
    this.onMouseUp = this.onMouseUp.bind(this);

    this.state = {
      x: this.props.initialX,
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
    if (e.button !== 0)
      return;
    e.stopPropagation();
    e.preventDefault();
    this.setState({
      dragging: true,
    });
  }

  onMouseUp(e) {
    e.stopPropagation();
    e.preventDefault();
    this.setState({
      dragging: false,
    });
  }

  onMouseMove(e) {
    if (!this.state.dragging)
      return;
    const x = e.pageX;
    const xMin = this.props.min || 0;
    const xMax = this.props.max || 9999999;
    if (x >= xMin && x <= xMax) {
      // TODO: maybe can set th
      this.setState({
        x: x
      });
      this.props.onPosChanged(x);
      //store.setSidebarDx(e.pageX);
    }
    e.stopPropagation();
    e.preventDefault();
  }

  render() {
    const style = {
      position: 'absolute',
      backgroundColor: '#377CE4',
      minHeight: '100%',
      width: 3,
      cursor: 'col-resize',
      zIndex: 3,
      left: this.state.x
    };

    return (
      <div
        style={style}
        onMouseDown={this.onMouseDown}
        onMouseMove={this.onMouseMove}
        onMouseUp={this.onMouseUp}>
      </div>
    );
  }
}
