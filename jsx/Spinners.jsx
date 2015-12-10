'use strict';

var React = require('react')

class SpinnerCircle extends React.Component {
	render() {
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
    )
	}
}

module.exports = {
  Circle: SpinnerCircle
};

