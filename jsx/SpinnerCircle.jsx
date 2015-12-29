import React from 'react';
import action from './action.js';

export default class SpinnerCircle extends React.Component {

  constructor(props, context) {
    super(props, context);
    this.handleToggleSpinner = this.handleToggleSpinner.bind(this);
    this.state = {
      visible: action.spinnerIsVisible()
    };
  }

  handleToggleSpinner(newVisibleState) {
    this.setState({visible: newVisibleState});
  }
  
  componentWillMount() {
    if (!this.props.forceVisible) {
      this.cidSpinner = action.onSpinner(this.handleToggleSpinner);
    }
  }

  componentWillUnmount() {
    if (!this.props.forceVisible) {
      action.offSpinner(this.cidSpinner);
    }
  }

  render() {
    const isVisible = this.props.forceVisible || this.state.visible;
    if (!isVisible) {
      return null;
    }
    return (
      <div style={this.props.style} className="circle-wrapper fade-in spinner">
        <div className="circle1 circle"></div>
        <div className="circle2 circle"></div>
        <div className="circle3 circle"></div>
        <div className="circle4 circle"></div>
        <div className="circle5 circle"></div>
        <div className="circle6 circle"></div>
        <div className="circle7 circle"></div>
        <div className="circle8 circle"></div>
        <div className="circle9 circle"></div>
        <div className="circle10 circle"></div>
        <div className="circle11 circle"></div>
        <div className="circle12 circle"></div>
      </div>
    );
  }
}
