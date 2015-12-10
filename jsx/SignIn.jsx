/* jshint -W097,-W117 */
'use strict';

var React = require('react');

class SignIn extends React.Component {
  render() {
    return (
      <div className="overlay-dialog overlay-dialog--signin " tabindex="-1">
        <h3 className="overlay-title">DB Hero</h3>
        <div className="overlay-content">Sign in to DB Hero or create an account</div>
        <div className="overlay-actions">
          <div className="buttonSet--vertical signin-auth-choices">
            <button className="button button--twitter" data-action="twitter-auth" data-action-source="nav_signup" title="Connect with Twitter" data-redirect="https://medium.com:443/">
              <span className="icon icon--twitter" style="display: none !important;"></span>
              <span className="button-label--twitter">Sign in with Twitter</span>
            </button>
            <button className="button button--facebook" data-action="facebook-auth" data-action-source="nav_signup" title="Connect with Facebook" data-redirect="https://medium.com:443/">
              <span className="icon icon--facebook" style="display: none !important;"></span>
              <span className="button-label--facebook">Sign in with Facebook</span>
            </button>
          </div>
        </div>
      </div>
    );
  }
}

module.exports = SignIn;
