/* jshint -W097,-W117 */
'use strict';

var React = require('react');

var action = require('./action.js');
var Modal = require('react-modal');

class QueryEditBar extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleSaveChanges = this.handleSaveChanges.bind(this);
    this.handleSQLPreview = this.handleSQLPreview.bind(this);

    // 1) Is there a way to move discard changes to here without using action?
    // 2) maybe move generateQuery from output to here?

    this.state = {
      modalIsOpen: false,
      modalText: "",
    };
  }

  closeModal() {
    this.setState({modalIsOpen: false});
  }

  handleModalCloseRequest() {
    // opportunity to validate something and keep the modal open even if it
    // requested to be closed
    this.closeModal();
  }


  handleSaveChanges() {
    console.log("handleSaveChanges ");

    // TODO: must support multiple queries for multiple rows changes
    var query = this.props.generateQuery();
    action.executeQuery(query);
  }

  handleSQLPreview() {
    console.log("handleSQLPreview");
    var query = this.props.generateQuery();

    this.setState({
      modalIsOpen: true,
      modalText: query,
    });
  }

  render() {
    var modalStyle = {
      content : {
        display               : 'block',
        overflow              : 'auto',
        top                   : '40%',
        left                  : '50%',
        right                 : 'none',
        bottom                : 'none',
        maxWidth              : '60%',
        transform             : 'translate(-50%, -50%)',
        fontSize              : '12px',
        background            : '#3B8686',
        color                 : '#fff',
      }
    };

    var appElement = document.getElementById('main');
    Modal.setAppElement(appElement);

    return (
      <div id="query_edit_bar">
        <button className="discard_changes" onClick={this.props.onHandleDiscardChanges}>Discard Changes</button>
        <div className="row_number">{this.props.numberOfRowsEdited} edited rows</div>
        <button className="sql_preview" onClick={this.handleSQLPreview.bind(this)}>SQL Preview</button>
        <button className="save_changes" onClick={this.handleSaveChanges.bind(this)}>Save Changes</button>

        <Modal
            isOpen={this.state.modalIsOpen}
            onRequestClose={this.closeModal.bind(this)}
            style={modalStyle} >
            {this.state.modalText}
        </Modal>
      </div>
    );
  }
}

module.exports = QueryEditBar;
