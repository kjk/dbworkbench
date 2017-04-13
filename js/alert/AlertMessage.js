import React from "react";
import ReactDOM from "react-dom";
import PropTypes from "prop-types";
import classnames from "classnames";

class AlertMessage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      closeButtonStyle: {},
    };
  }

  componentDidMount() {
    this.domNode = ReactDOM.findDOMNode(this);
    this.setState({
      closeButtonStyle: {
        height: this.domNode.offsetHeight + "px",
        lineHeight: this.domNode.offsetHeight + "px",
        backgroundColor: this.props.style.closeButton.bg,
      },
    });

    if (this.props.time > 0) {
      this._countdown();
    }
  }

  _removeSelf() {
    reactAlertEvents.emit("ALERT.REMOVE", this);
  }

  _countdown() {
    setTimeout(() => {
      this._removeSelf();
    }, this.props.time);
  }

  _showIcon() {
    let icon = "";
    if (this.props.icon) {
      icon = this.props.icon;
    } else {
      icon = <div className={this.props.type + "-icon"} />;
    }

    return icon;
  }

  _handleCloseClick() {
    this._removeSelf();
  }

  render() {
    return (
      <div
        style={this.props.style.alert}
        className={classnames("alert", this.props.type)}
      >
        <div className="content icon">
          {this._showIcon.bind(this)()}
        </div>
        <div className="content message">
          {this.props.message}
        </div>
        <div
          onClick={this._handleCloseClick.bind(this)}
          style={this.state.closeButtonStyle}
          className="content close"
        >
          <div className={this.props.closeIconClass} />
        </div>
      </div>
    );
  }
}

AlertMessage.defaultProps = {
  reactKey: "",
  icon: "",
  message: "",
  type: "info",
};

AlertMessage.propTypes = {
  type: PropTypes.oneOf(["info", "success", "error"]),
  time: PropTypes.number
};

export default AlertMessage;
