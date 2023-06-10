package com.util.auth;

import com.common.exception.AccessException;
import com.common.policy.AccessVersionClient;
import com.common.policy.AccessVersionPolicy;
import com.common.policy.BandwidthPolicy;
import com.util.db.AccessVersionPolicyDao;
import com.util.message.I18nMessageUtil;
import com.util.misc.StrUtil;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;

import static com.common.constant.ResponseCode.RSP_CLIENT_ACCESS_LIMITED;
import static com.common.constant.ResponseCode.RSP_CLIENT_NEED_UPDATE;

/**
 * @ClassName AccessVersionUtil
 * @Author xxxx
 * @Date 2023.5.17 19:09
 * @since JDK 1.8
 */
@Slf4j
@Component
public class AccessVersionUtil {
    @Autowired
    private AccessVersionPolicyDao accessVersionPolicyDao;

    @Autowired
    private StrUtil strUtil;

    @Autowired
    private I18nMessageUtil i18nMessageUtil;

//    @Autowired
//    private BandwidthPolicy bandwidthPolicy;

    public void checkClientLoginStrategy(AccessVersionClient client) {
        List<AccessVersionPolicy> allAccessVersionPolicy = accessVersionPolicyDao.getAllEnableAccessVersionPolicy();

        if (CollectionUtils.isEmpty(allAccessVersionPolicy)) {
            // 没用策略，无限制，正常返回。
            log.error("no limit");
            return;
        }

        for (AccessVersionPolicy policy : allAccessVersionPolicy) {
            if (isPolicyLimit(client, policy) == true) {
                // 除版本外的参数满足策略，则继续检查版本
                checkVersion(client, policy);
            } else {
                continue;
            }
        }
    }

    public boolean isPolicyLimit(AccessVersionClient client, AccessVersionPolicy policy) {
        if ((policy.getTerminalMode() == -1 || client.getClientType() == policy.getTerminalMode())
                && (policy.getTerminalSysType() == -1 || client.getOsType() == policy.getTerminalSysType())
                && (policy.getNetworkType() == -1 || client.getNetType() == policy.getNetworkType())
                && (policy.getHardwareType() == -1 || client.getHardware() == policy.getHardwareType())) {
            return true;
        } else {
            return false;
        }
    }

    private void checkVersion(AccessVersionClient client, AccessVersionPolicy policy) {
        if (policy.getMinimumVersionId() == -1) {
            throw new AccessException(RSP_CLIENT_ACCESS_LIMITED, i18nMessageUtil.getErrorMessage(RSP_CLIENT_ACCESS_LIMITED));
        } else {
            //客户端可能会携带多个版本号，以逗号分隔
            checkClientVersionList(client, policy);
        }
    }

    private void checkClientVersionList(AccessVersionClient client, AccessVersionPolicy policy) {
        List<String> clientVersionList = strUtil.splitList(client.getClientVersion(), ",");

        for (String clientVersion : clientVersionList) {
            oneVersionChecker(clientVersion, policy.getMinimumVersionName());
        }
    }

    private void oneVersionChecker(String clientVersion, String minimumVersion) {
        String clientVer = StringUtils.replaceChars(clientVersion, "V_", "");
        String minimumVer = StringUtils.replaceChars(minimumVersion, "V_", "");

        if (StringUtils.equalsIgnoreCase(clientVer, minimumVer)) {
            // 版本相同
            log.error("client version : {}, minimum version : {}", clientVer, minimumVer);
            throw new AccessException(RSP_CLIENT_NEED_UPDATE, i18nMessageUtil.getErrorMessage(RSP_CLIENT_NEED_UPDATE));
        }

        List<String> clientVerList = strUtil.splitList(clientVer, ".");
        List<String> minimumVerList = strUtil.splitList(minimumVer, ".");

        for (int clientVerIdx = 0, minimumVerIdx = 0;
             clientVerIdx < clientVerList.size() && minimumVerIdx < minimumVerList.size();
             clientVerIdx++, minimumVerIdx++) {
            String clientVerMinor = clientVerList.get(clientVerIdx);
            String minimumVerMinor = minimumVerList.get(minimumVerIdx);

            int result = StringUtils.compare(clientVerMinor, minimumVerMinor);
            if (result < 0) {
                // 客户端版本小，需要升级
                log.error("client version : {}, minimum version : {}", clientVer, minimumVer);
                throw new AccessException(RSP_CLIENT_NEED_UPDATE, i18nMessageUtil.getErrorMessage(RSP_CLIENT_NEED_UPDATE));
            } else if (result > 0) {
                // 客户端版本小
                return;
            }
        }
    }
}