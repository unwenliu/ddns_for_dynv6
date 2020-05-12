# coding=utf8
import sys
from json import load, dump
import logging
from argparse import ArgumentParser, RawTextHelpFormatter
import subprocess
import re
import requests

__description__ = "自动更新dynv6的DNS记录指向本地ip,支持ipv4和ipv6"
__version__ = '%(prog)s 1.0'

# 输入log到文件
logging.basicConfig(filename='ddns_for_dynv6.log', level=logging.INFO, filemode='a',
                    format='%(asctime)s - %(levelname)s: %(message)s')


def get_config(path="config.json"):
    """
    创建配置文件
    :param path:默认当前目录下
    :return: json结构的config
    """
    try:
        with open(path) as config_file:
            config = load(config_file)
    except IOError:
        logging.error("配置文件{}不存在".format(path))
        with open(path, 'w') as config_file:
            config = {
                "hostname": "your hostname",
                "token": "your token",
                "ipv4": True,
                "ipv6": True
            }
            dump(config, config_file, indent=2, sort_keys=True)
        sys.exit("已经根据模板生成新的配置文件{}".format(path))
    except:
        sys.exit("不能够加载配置文件:{}".format(path))
    return config


def get_ip():
    """
    获取ip
    :return:返回字典形式ip地址
    """
    child = subprocess.Popen('ipconfig', shell=True, stdout=subprocess.PIPE)
    out = str(child.communicate())

    ipv4_pattern = '(?:[0-9]{1,3}\.){3}[0-9]{1,3}'
    ipv6_pattern = '(([a-f0-9]{1,4}:){7}[a-f0-9]{1,4})'

    all_ipv4_address = re.findall(ipv4_pattern, out)
    all_ipv6_address = re.findall(ipv6_pattern, out)

    ipv4_address = all_ipv4_address[0]
    ipv6_address = all_ipv6_address[0][0]
    return {"ipv4": ipv4_address, "ipv6": ipv6_address}


def res(**config):
    """
    向dyn6发送请求
    :return:
    """
    url = 'https://dynv6.com/api/update'
    r = requests.get(url, params=config)
    return r


def ipv4_log(response):
    if response.status_code != 200:
        logging.error("更新ipv4地址失败，原因是dynv6返回{}，{}的异常响应".format(r.status_code, r.text))
    else:
        logging.info("更新ipv4地址成功，dynv6返回{}的正常响应".format(r.text))


def ipv6_log(response):
    if response.status_code != 200:
        logging.error("更新ipv6地址失败，原因是dynv6返回{}，{}的异常响应".format(r.status_code, r.text))
    else:
        logging.info("更新ipv6地址成功，dynv6返回{}的正常响应".format(r.text))


def ipv4_ipv6_log(response):
    if response.status_code != 200:
        logging.error("更新ipv4地址和ipv6地址失败，原因是dynv6返回{}，{}的异常响应".format(r.status_code, r.text))
    else:
        logging.info("更新ipv4地址和ipv6地址成功，dynv6返回{}的正常响应".format(r.text))


def parse_args():
    """
    进行参数解析
    :return:
    """
    parser = ArgumentParser(description=__description__, formatter_class=RawTextHelpFormatter)
    use_config_group = parser.add_argument_group(title='使用配置文件运行')
    use_args_group = parser.add_argument_group(title='使用参数运行')
    use_config_group.add_argument('-c', '--config', help="指定配置文件运行 [配置文件路径]")
    use_args_group.add_argument('-hostname', help="要更新的域名")
    use_args_group.add_argument('-token', help="dynv6里域名所绑定的token")
    use_args_group.add_argument('-4', '--ipv4', action='store_true', help="更新ipv4地址", dest='ipv4')
    use_args_group.add_argument('-6', '--ipv6', action='store_true', help="更新ipv6地址", dest='ipv6')
    parser.add_argument('-v', '--version', action='version', version=__version__, help='显示版本信息')
    args = parser.parse_args()
    return args


if __name__ == '__main__':
    ip = get_ip()
    args = parse_args()
    # 使用参数运行
    if args.ipv4 and args.ipv6 is False and args.hostname is not None and args.token is not None:
        r = res(hostname=args.hostname, token=args.token, ipv4=ip["ipv4"])
        ipv4_log(r)
    if args.ipv6 and args.ipv4 is False and args.hostname is not None and args.token is not None:
        r = res(hostname=args.hostname, token=args.token, ipv6=ip["ipv6"])
        ipv6_log(r)
    if args.ipv4 and args.ipv6 and args.hostname is not None and args.token is not None:
        r = res(hostname=args.hostname, token=args.token, ipv4=ip["ipv4"], ipv6=ip["ipv6"])
        ipv4_ipv6_log(r)
    # 使用配置文件运行
    if args.config is not None:
        config = get_config(args.config)
        if config["ipv4"] and config["ipv6"] is False:
            r = res(hostname=config["hostname"], token=config["token"], ipv4=ip["ipv4"])
            ipv4_log(r)
        if config["ipv6"] and config["ipv4"] is False:
            r = res(hostname=config["hostname"], token=config["token"], ipv6=ip["ipv6"])
            ipv6_log(r)
        elif config["ipv4"] and config["ipv6"]:
            r = res(hostname=config["hostname"], token=config["token"], ipv4=ip["ipv4"], ipv6=ip["ipv6"])
            ipv4_ipv6_log(r)
