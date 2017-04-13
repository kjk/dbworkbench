import React from "react";
import ReactDOM from "react-dom";
import PropTypes from "prop-types";

export default class DragBarVert extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.handleMouseDown = this.handleMouseDown.bind(this);
    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseUp = this.handleMouseUp.bind(this);

    this.x = this.props.initialX;

    this.state = {
      dragging: false,
    };
  }

  componentDidUpdate(props, prevState) {
    const dragEnter = !prevState.dragging && this.state.dragging;
    const dragLeave = prevState.dragging && !this.state.dragging;
    if (dragEnter) {
      document.addEventListener("mousemove", this.handleMouseMove);
      document.addEventListener("mouseup", this.handleMouseUp);
    } else if (dragLeave) {
      document.removeEventListener("mousemove", this.handleMouseMove);
      document.removeEventListener("mouseup", this.handleMouseUp);
    }
  }

  handleMouseDown(e) {
    // only left mouse button
    if (e.button !== 0) return;
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

    const x = e.pageX;
    const xMin = this.props.min || 0;
    const xMax = this.props.max || 9999999;
    if (x >= xMin && x <= xMax) {
      this.x = x;
      const el = ReactDOM.findDOMNode(this);
      el.style.left = x + "px";
      this.props.onPosChanged(x);
    }
    e.stopPropagation();
    e.preventDefault();
  }

  render() {
    const style = {
      position: "absolute",
      backgroundColor: "#377CE4",
      minHeight: "100%",
      width: 3,
      cursor: "col-resize",
      zIndex: 3,
      left: this.x,
    };

    return <div style={style} onMouseDown={this.handleMouseDown} />;
  }
}

DragBarVert.propTypes = {
  onPosChanged: PropTypes.func.isRequired,
  initialX: PropTypes.number.isRequired,
  min: PropTypes.number.isRequired,
  max: PropTypes.number.isRequired,
};
