import 'babel-polyfill';
import React from 'react';
import ReactDOM from 'react-dom';
import * as store from './store.js';

// TODO: add propTypes
// http://www.newmediacampaigns.com/blog/refactoring-react-components-to-es6-classes
export default class DragBarVert extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.handleMouseDown = this.handleMouseDown.bind(this);
    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseUp = this.handleMouseUp.bind(this);

    this.state = {
      x: this.props.initialX,
      dragging: false,
    };
  }

  componentDidUpdate(props, state) {
    if (this.state.dragging && !state.dragging) {
      document.addEventListener('mousemove', this.handleMouseMove);
      document.addEventListener('mouseup', this.handleMouseUp);
      //console.log("DragBarVert adding event handlers");
    } else if (!this.state.dragging && state.dragging) {
      document.removeEventListener('mousemove', this.handleMouseMove);
      document.removeEventListener('mouseup', this.handleMouseUp);
      //console.log("DragBarVert removing event handlers");
    }
  }

  handleMouseDown(e) {
    // only left mouse button
    if (e.button !== 0)
      return;
    e.stopPropagation();
    e.preventDefault();
    this.setState({
      dragging: true,
    });
  }

  handleMouseUp(e) {
    e.stopPropagation();
    e.preventDefault();
    this.setState({
      dragging: false,
    });
  }

  handleMouseMove(e) {
    if (!this.state.dragging) {
      //console.log("DragBarVert not dragging");
      return;
    }

    const x = e.pageX;
    const xMin = this.props.min || 0;
    const xMax = this.props.max || 9999999;
    if (x >= xMin && x <= xMax) {
      // TODO: maybe can set the element.style.width directly?
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
        onMouseDown={this.handleMouseDown}
      ></div>
    );

    /*
    return (
      <div
        style={style}
        onMouseDown={this.handleMouseDown}
        onMouseMove={this.handleMouseMove}
        onMouseUp={this.handleMouseUp}>
      </div>
    );*/

  }
}
