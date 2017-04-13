import React from 'react';
import PropTypes from 'prop-types';
import * as store from './store.js';

export default class SpinnerCircle extends React.Component {

  constructor(props, context) {
    super(props, context);
    this.handleToggleSpinner = this.handleToggleSpinner.bind(this);
    this.state = {
      visible: store.spinnerIsVisible()
    };
  }

  componentWillMount() {
    if (!this.props.forceVisible) {
      store.onSpinner(this.handleToggleSpinner, this);
    }
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  handleToggleSpinner(newVisibleState) {
    this.setState({
      visible: newVisibleState
    });
  }

  render() {
    const isVisible = this.props.forceVisible || this.state.visible;
    if (!isVisible) {
      return null;
    }
    return (
      <div style={ this.props.style } className='circle-wrapper fade-in spinner'>
        <div className='circle1 circle' />
        <div className='circle2 circle' />
        <div className='circle3 circle' />
        <div className='circle4 circle' />
        <div className='circle5 circle' />
        <div className='circle6 circle' />
        <div className='circle7 circle' />
        <div className='circle8 circle' />
        <div className='circle9 circle' />
        <div className='circle10 circle' />
        <div className='circle11 circle' />
        <div className='circle12 circle' />
      </div>
      );
  }
}

SpinnerCircle.propTypes = {
  forceVisible: PropTypes.bool,
  style: PropTypes.object
};
