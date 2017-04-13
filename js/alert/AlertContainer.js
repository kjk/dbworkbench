import React from "react";
import PropTypes from "prop-types";
import EventEmitter from "events";
import ReactCSSTransitionGroup from "react-addons-css-transition-group";
import AlertMessage from "./AlertMessage";

class AlertContainer extends React.Component {
  constructor(props) {
    super(props);
    global.reactAlertEvents = new EventEmitter();
    this.state = {
      alerts: [],
    };
    this.style = this._setStyle();
    this.theme = this._setTheme();
    this._eventListners();
  }

  success(message, options = {}) {
    options.type = "success";
    this.show(message, options);
  }

  error(message, options = {}) {
    options.type = "error";
    this.show(message, options);
  }

  info(message, options = {}) {
    options.type = "info";
    this.show(message, options);
  }

  show(message, options = {}) {
    let alert = {};
    alert.message = message;
    alert = Object.assign(alert, options);
    this.setState({ alerts: this._addAlert(alert) });
  }

  removeAll() {
    this.setState({ alerts: [] });
  }

  _addAlert(alert) {
    alert.uniqueKey = this._genUniqueKey();
    alert.style = this.theme;
    if (!alert.hasOwnProperty("time")) {
      alert.time = this.props.time;
    }
    alert.closeIconClass = "close-" + this.props.theme;
    this.state.alerts.push(alert);
    return this.state.alerts;
  }

  _genUniqueKey() {
    return (
      new Date().getTime().toString() + Math.random().toString(36).substr(2, 5)
    );
  }

  _removeAlert(alert) {
    return this.state.alerts.filter(a => {
      return a.uniqueKey != alert.props.uniqueKey;
    });
  }

  _eventListners() {
    reactAlertEvents.on("ALERT.REMOVE", alert => {
      this.setState({ alerts: this._removeAlert(alert) });
    });
  }

  _setStyle() {
    let position = {};
    switch (this.props.position) {
      case "top left":
        position = {
          top: 0,
          right: "auto",
          bottom: "auto",
          left: 0,
        };
        break;
      case "top right":
        position = {
          top: 0,
          right: 0,
          bottom: "auto",
          left: "auto",
        };
        break;
      case "bottom left":
        position = {
          top: "auto",
          right: "auto",
          bottom: 0,
          left: 0,
        };
        break;
      default:
        position = {
          top: "auto",
          right: 0,
          bottom: 0,
          left: "auto",
        };
        break;
    }

    return {
      margin: this.props.offset + "px",
      top: position.top,
      right: position.right,
      bottom: position.bottom,
      left: position.left,
    };
  }
  /**
   * Set the style of the alert based on the chosen theme
   */
  _setTheme() {
    let theme = {};
    switch (this.props.theme) {
      case "light":
        theme = {
          alert: {
            backgroundColor: "#fafafa",
            color: "#333",
          },
          closeButton: {
            bg: "#fafafa",
          },
        };
        break;
      default:
        theme = {
          alert: {
            backgroundColor: "#333",
            color: "#fff",
          },
          closeButton: {
            bg: "#444",
          },
        };
        break;
    }

    return theme;
  }

  componentDidUpdate() {
    this.style = this._setStyle();
    this.theme = this._setTheme();
  }

  render() {
    return (
      <div style={this.style} className="react-alerts">
        <ReactCSSTransitionGroup
          transitionName={this.props.transition}
          transitionEnterTimeout={250}
          transitionLeaveTimeout={250}
        >
          {this.state.alerts.map((alert, index) => {
            return <AlertMessage key={alert.uniqueKey} {...alert} />;
          })}
        </ReactCSSTransitionGroup>
      </div>
    );
  }
}

AlertContainer.defaultProps = {
  offset: 14,
  position: "bottom left",
  theme: "dark",
  time: 5000,
  transition: "scale",
};

AlertContainer.propTypes = {
  offset: PropTypes.number,
  position: PropTypes.oneOf([
    "bottom left",
    "bottom right",
    "top right",
    "top left",
  ]),
  theme: PropTypes.oneOf(["dark", "light"]),
  time: PropTypes.number,
  transition: PropTypes.oneOf(["scale", "fade"]),
};

export default AlertContainer;
