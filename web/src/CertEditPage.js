import React from "react";
import {Button, Card, Col, Input, Row, Select} from "antd";
import * as CertBackend from "./backend/CertBackend";
import * as Setting from "./Setting";
import i18next from "i18next";
import copy from "copy-to-clipboard";
import FileSaver from "file-saver";

const {Option} = Select;
const {TextArea} = Input;

class CertEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      certName: props.match.params.certName,
      cert: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getCert();
  }

  getCert() {
    CertBackend.getCert(this.props.account.name, this.state.certName)
      .then((cert) => {
        this.setState({
          cert: cert,
        });
      });
  }

  parseCertField(key, value) {
    if (["score"].includes(key)) {
      value = Setting.myParseInt(value);
    }
    return value;
  }

  updateCertField(key, value) {
    value = this.parseCertField(key, value);

    const cert = this.state.cert;
    cert[key] = value;
    this.setState({
      cert: cert,
    });
  }

  renderCert() {
    const editorWidth = Setting.isMobile() ? 22 : 9;
    return (
      <Card size="small" title={
        <div>
          {i18next.t("cert:Edit Cert")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="primary" onClick={this.submitCertEdit.bind(this)}>{i18next.t("general:Save")}</Button>
        </div>
      } style={{marginLeft: "5px"}} type="inner">
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Name")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.cert.name} onChange={e => {
              this.updateCertField("name", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Display name")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.cert.displayName} onChange={e => {
              this.updateCertField("displayName", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("cert:Type")}:
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: "100%"}} value={this.state.cert.type} onChange={(value => {this.updateCertField("type", value);})}>
              {
                [
                  {id: "SSL", name: "SSL"},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("cert:Crypto algorithm")}:
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: "100%"}} value={this.state.cert.cryptoAlgorithm} onChange={(value => {this.updateCertField("cryptoAlgorithm", value);})}>
              {
                [
                  {id: "RSA", name: "RSA"},
                  {id: "ECC", name: "ECC"},
                ].map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("cert:Expire time")}:
          </Col>
          <Col span={22} >
            <Input disabled={true} value={Setting.getFormattedDate(this.state.cert.expireTime)} onChange={e => {
              this.updateCertField("expireTime", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("cert:Certificate")}:
          </Col>
          <Col span={editorWidth} >
            <Button style={{marginRight: "10px", marginBottom: "10px"}} onClick={() => {
              copy(this.state.cert.certificate);
              Setting.showMessage("success", i18next.t("cert:Certificate copied to clipboard successfully"));
            }}
            >
              {i18next.t("cert:Copy certificate")}
            </Button>
            <Button type="primary" onClick={() => {
              const blob = new Blob([this.state.cert.certificate], {type: "text/plain;charset=utf-8"});
              FileSaver.saveAs(blob, "token_jwt_key.pem");
            }}
            >
              {i18next.t("cert:Download certificate")}
            </Button>
            <TextArea autoSize={{minRows: 30, maxRows: 30}} value={this.state.cert.certificate} onChange={e => {
              this.updateCertField("certificate", e.target.value);
            }} />
          </Col>
          <Col span={1} />
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {i18next.t("cert:Private key")}:
          </Col>
          <Col span={editorWidth} >
            <Button style={{marginRight: "10px", marginBottom: "10px"}} onClick={() => {
              copy(this.state.cert.privateKey);
              Setting.showMessage("success", i18next.t("cert:Private key copied to clipboard successfully"));
            }}
            >
              {i18next.t("cert:Copy private key")}
            </Button>
            <Button type="primary" onClick={() => {
              const blob = new Blob([this.state.cert.privateKey], {type: "text/plain;charset=utf-8"});
              FileSaver.saveAs(blob, "token_jwt_key.key");
            }}
            >
              {i18next.t("cert:Download private key")}
            </Button>
            <TextArea autoSize={{minRows: 30, maxRows: 30}} value={this.state.cert.privateKey} onChange={e => {
              this.updateCertField("privateKey", e.target.value);
            }} />
          </Col>
        </Row>
      </Card>
    );
  }

  submitCertEdit() {
    const cert = Setting.deepCopy(this.state.cert);
    CertBackend.updateCert(this.state.cert.owner, this.state.certName, cert)
      .then((res) => {
        if (res) {
          Setting.showMessage("success", "Successfully saved");
          this.setState({
            certName: this.state.cert.name,
          });
          this.props.history.push(`/certs/${this.state.cert.name}`);
        } else {
          Setting.showMessage("error", "failed to save: server side failure");
          this.updateCertField("name", this.state.certName);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `failed to save: ${error}`);
      });
  }

  render() {
    return (
      <div>
        <Row style={{width: "100%"}}>
          <Col span={1}>
          </Col>
          <Col span={22}>
            {
              this.state.cert !== null ? this.renderCert() : null
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
        <Row style={{margin: 10}}>
          <Col span={2}>
          </Col>
          <Col span={18}>
            <Button type="primary" size="large" onClick={this.submitCertEdit.bind(this)}>{i18next.t("general:Save")}</Button>
          </Col>
        </Row>
      </div>
    );
  }
}

export default CertEditPage;
