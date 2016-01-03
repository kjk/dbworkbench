import 'babel-polyfill';
import React from 'react';
import ReactDOM from 'react-dom';

// TODO: add propTypes
// http://www.newmediacampaigns.com/blog/refactoring-react-components-to-es6-classes
export default class DragBarHoriz extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.handleMouseDown = this.handleMouseDown.bind(this);
    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseUp = this.handleMouseUp.bind(this);

    this.y = this.props.initialY;
    this.state = {
      dragging: false,
    };
  }

  componentDidUpdate(props, prevState) {
    const dragEnter = !prevState.dragging && this.state.dragging;
    const dragLeave = prevState.dragging && !this.state.dragging;
    if (dragEnter) {
      document.addEventListener('mousemove', this.handleMouseMove);
      document.addEventListener('mouseup', this.handleMouseUp);
    } else if (dragLeave) {
      document.removeEventListener('mousemove', this.handleMouseMove);
      document.removeEventListener('mouseup', this.handleMouseUp);
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
      return;
    }

    const y = e.pageY;
    const yMin = this.props.min || 0;
    const yMax = this.props.max || 9999999;
    if (y >= yMin && y <= yMax) {
      this.y = y;
      const el = ReactDOM.findDOMNode(this);
      el.style.top = y + 'px';
      this.props.onPosChanged(y);
    }
    e.stopPropagation();
    e.preventDefault();
  }

  render() {
    const style = {
      position: 'absolute',
      backgroundColor: '#377CE4',
      minWidth: '100%',
      height: 3,
      cursor: 'row-resize',
      zIndex: 0,
      top: this.y
    };

    return (
      <div style={ style } onMouseDown={ this.handleMouseDown }>
      </div>
      );
  }
}
